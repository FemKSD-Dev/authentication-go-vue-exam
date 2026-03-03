package bootstrap

import (
	"database/sql"
	"fmt"
	"log"

	"authentication-project-exam/migrations"

	"github.com/pressly/goose/v3"
)

func RunMigrations(db *sql.DB) error {
	goose.SetBaseFS(migrations.Files)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("could not run migrations: %w", err)
	}
	log.Println("Migrations completed successfully")
	return nil
}
