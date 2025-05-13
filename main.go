package main

import (
	"net/http"
	"time"

	"github.com/edwinboon/workout-tracking-api/internal/app"
)

func main() {
	// Initialize the application
	app, err := app.NewApplication()

	if err != nil {
		panic(err) // Self destruct if we can't start the application
	}

	// Start the application
	app.Logger.Println("Starting application...")

	// Setup server
	server := &http.Server{
		Addr:         ":8080",
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Running server

	err = server.ListenAndServe()

	if err != nil {
		app.Logger.Fatal(err)
	}
}
