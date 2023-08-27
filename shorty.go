package main

import (
	"errors"
	"log"
	"net/http"
	"shorty/internal/api"
	"shorty/internal/config"
	"shorty/internal/shorty"
	"shorty/internal/storage/sqlite"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sqlx.Open("sqlite3", config.C().DBpath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	storage := sqlite.NewStorage(db)
	shorty := shorty.New(storage.Shortenings)
	api := api.New(shorty)

	go func() {
		server := http.Server{
			Addr:    config.C().Address,
			Handler: api,
		}

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error starting server: %v", err)
		}
	}()

	time.Sleep(100 * time.Second)
}
