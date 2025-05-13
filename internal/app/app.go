package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/edwinboon/workout-tracking-api/internal/api"
	"github.com/edwinboon/workout-tracking-api/internal/store"
	"github.com/edwinboon/workout-tracking-api/migrations"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	DB             *sql.DB
}

func NewApplication() (*Application, error) {

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// stores
	pgDB, err := store.Open()

	if err != nil {
		return nil, err
	}

	// migrations
	err = store.MigrateFS(pgDB, migrations.FS, ".")

	if err != nil {
		panic(err) // if database is not working just self-destruct
	}

	// handlers
	workoutHandler := api.NewWorkoutHandler()

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		DB:             pgDB,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
