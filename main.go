package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/edwinboon/workout-tracking-api/internal/app"
	"github.com/edwinboon/workout-tracking-api/internal/routes"
)

func main() {
	// Make port configurable
	var port int
	flag.IntVar(&port, "port", 8080, "Port to run the server on")
	flag.Parse()

	// Initialize the application
	app, err := app.NewApplication()

	if err != nil {
		panic(err) // Self destruct if we can't start the application
	}

	// Setup server
	r := routes.SetupRoutes(app)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("Starting server on port %d\n", port)

	// Running server
	err = server.ListenAndServe()

	if err != nil {
		app.Logger.Fatal(err)
	}
}
