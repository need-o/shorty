package migrate

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func RunForSqlite3(db *sql.DB, source string) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		fmt.Println(1)
		return err
	}

	migrate, err := migrate.NewWithDatabaseInstance(source, "sqlite3", driver)
	if err != nil {
		return err
	}

	return migrate.Up()
}
