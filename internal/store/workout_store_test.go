package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")

	if err != nil {
		t.Fatalf("opening test db: %v", err)
	}

	// run migrations for our test db
	err = Migrate(db, "../../migrations/")

	if err != nil {
		t.Fatalf("migrating test db: %v", err)
	}

	// every time we run a test, we want to start with a clean slate
	_, err = db.Exec("TRUNCATE TABLE workouts, workout_entries CASCADE")

	if err != nil {
		t.Fatalf("truncating test db: %v", err)
	}

	return db

}
