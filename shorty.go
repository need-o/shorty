package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"shorty/internal/api"
	"shorty/internal/config"
	"shorty/internal/shorty"
	"shorty/internal/storage/sqlite"
	"syscall"
	"time"

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

	storage := sqlite.NewStorage(db)
	shorty := shorty.New(storage.Shorty)
	api := api.New(shorty)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

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

	api.AddCloser(
		func(context.Context) error {
			return db.Close()
		},
	)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := api.Shutdown(ctx); err != nil {
		log.Fatalf("error stopping server: %v", err)
	}
}
