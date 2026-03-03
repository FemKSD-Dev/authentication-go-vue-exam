package main

import (
	"authentication-project-exam/internal/bootstrap"
	"log"
)

func main() {
	app, err := bootstrap.NewApp()
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}

	log.Println("Application started successfully")
}
