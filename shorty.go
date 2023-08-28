package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/need-o/shorty/internal/api"
	"github.com/need-o/shorty/internal/config"
	"github.com/need-o/shorty/internal/migrate"
	"github.com/need-o/shorty/internal/shorty"
	"github.com/need-o/shorty/internal/storage/sqlite"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

func main() {
	db, err := sqlx.Open("sqlite3", config.C().DBpath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if _, err := os.Stat(config.C().DBpath); errors.Is(err, os.ErrNotExist) {
		err := migrate.RunForSqlite3(db.DB, config.C().MigrationsSource)
		if err != nil {
			log.Fatal(err)
		}
	}

	storage := sqlite.NewStorage(db)
	shorty := shorty.New(storage.Shorty)
	api := api.New(shorty)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	api.AddCloser(
		func(context.Context) error {
			return db.Close()
		},
	)

	go func() {
		server := http.Server{
			Addr:    config.C().Address,
			Handler: api,
		}

		log.Infof("server listening in %v", server.Addr)

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error starting server: %v", err)
		}
	}()

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := api.Shutdown(ctx); err != nil {
		log.Fatalf("error stopping server: %v", err)
	}
}
