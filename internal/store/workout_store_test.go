package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresWorkoutStore(db)

	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Title:           "Push day",
				Description:     "Chest and triceps workout",
				DurationMinutes: 60,
				CaloriesBurned:  300,
				Entries: []WorkoutEntry{
					{
						ExerciseName:    "Bench Press",
						Sets:            3,
						Reps:            IntPtr(10),
						DurationSeconds: nil,
						Weight:          FloatPtr(72.5),
						Notes:           "Warmup sets",
						OrderIndex:      1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "workout with invalid entries",
			workout: &Workout{
				Title:           "Full body workout",
				Description:     "A mix of exercises",
				DurationMinutes: 90,
				CaloriesBurned:  500,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Plank",
						Sets:         3,
						Reps:         IntPtr(60),
						Notes:        "Keep in form",
						OrderIndex:   1,
					},
					{
						ExerciseName:    "Squats",
						Sets:            4,
						Reps:            IntPtr(15),
						DurationSeconds: IntPtr(30),
						Weight:          FloatPtr(100.0),
						Notes:           "Full depth",
						OrderIndex:      2,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(tt.workout)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.DurationMinutes, createdWorkout.DurationMinutes)
			assert.Equal(t, tt.workout.CaloriesBurned, createdWorkout.CaloriesBurned)

			retrievedWorkout, err := store.GetWorkoutByID(int64(createdWorkout.ID))

			require.NoError(t, err)

			assert.Equal(t, tt.workout.Title, retrievedWorkout.Title)
			assert.Equal(t, tt.workout.Description, retrievedWorkout.Description)
			assert.Equal(t, tt.workout.DurationMinutes, retrievedWorkout.DurationMinutes)
			assert.Equal(t, tt.workout.CaloriesBurned, retrievedWorkout.CaloriesBurned)
			assert.Equal(t, len(tt.workout.Entries), len(retrievedWorkout.Entries))

			// check entries
			for i := range retrievedWorkout.Entries {
				assert.Equal(t, tt.workout.Entries[i].ExerciseName, retrievedWorkout.Entries[i].ExerciseName)
				assert.Equal(t, tt.workout.Entries[i].Sets, retrievedWorkout.Entries[i].Sets)
				assert.Equal(t, tt.workout.Entries[i].Reps, retrievedWorkout.Entries[i].Reps)
				assert.Equal(t, tt.workout.Entries[i].DurationSeconds, retrievedWorkout.Entries[i].DurationSeconds)
				assert.Equal(t, tt.workout.Entries[i].Weight, retrievedWorkout.Entries[i].Weight)
				assert.Equal(t, tt.workout.Entries[i].Notes, retrievedWorkout.Entries[i].Notes)
				assert.Equal(t, tt.workout.Entries[i].OrderIndex, retrievedWorkout.Entries[i].OrderIndex)
			}

		})
	}
}

func IntPtr(i int) *int {
	return &i
}

func FloatPtr(i float64) *float64 {
	return &i
}
