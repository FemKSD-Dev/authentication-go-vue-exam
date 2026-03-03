package main

import (
	"authentication-project-exam/internal/bootstrap"
	"authentication-project-exam/internal/config"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	appLogger, err := bootstrap.NewLogger()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	db, err := bootstrap.NewDatabase(cfg, appLogger)
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	if err := bootstrap.RunMigrations(sqlDB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
}
