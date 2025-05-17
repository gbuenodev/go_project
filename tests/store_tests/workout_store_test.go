package store_test

import (
	"database/sql"
	"testing"

	"github.com/gbuenodev/goProject/internal/store"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	dbConfig := store.DBConfig{
		Provider: "Postgres",
		Driver:   "pgx",
		User:     "postgres",
		Password: "postgres",
		DBName:   "postgres",
		Host:     "localhost",
		Port:     5555,
		SSL:      "disable",
	}

	DBConn, err := store.Open(&dbConfig)
	if err != nil {
		t.Fatalf("Failed to open database connection: %v", err)
	}

	err = store.Migrate(DBConn, "../../migrations")
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	_, err = DBConn.Exec("TRUNCATE TABLE users, workouts, workout_entries RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("Failed to truncate tables: %v", err)
	}

	return DBConn
}

func TestCreateWorkout(t *testing.T) {
	DBConn := setupTestDB(t)

	workoutStore := store.NewPostgresWorkoutStore(DBConn)
	userStore := store.NewPostgresUserStore(DBConn)

	testUser := &store.User{
		Username: "Test_User",
		Email:    "test@email.com",
	}

	err := testUser.PasswordHash.Set("Sup3rSecr3tPass#!")
	require.NoError(t, err)

	err = userStore.CreateUser(testUser)
	require.NoError(t, err)

	tests := []struct {
		name    string
		workout *store.Workout
		wantErr bool
	}{
		{
			name: "Valid workout",
			workout: &store.Workout{
				UserID:          testUser.ID,
				Title:           "Push Day",
				Description:     "A workout focused on push exercises.",
				DurationMinutes: 60,
				CaloriesBurned:  500,
				Entries: []store.WorkoutEntry{
					{
						ExerciseName: "Bench Press",
						Reps:         IntPtr(10),
						Sets:         3,
						Weight:       FloatPtr(100),
						Notes:        "Felt strong today",
						OrderIndex:   1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid entries",
			workout: &store.Workout{
				UserID:          testUser.ID,
				Title:           "Leg Day",
				Description:     "A workout focused on leg exercises.",
				DurationMinutes: 45,
				CaloriesBurned:  400,
				Entries: []store.WorkoutEntry{
					{
						ExerciseName: "Squats",
						Reps:         IntPtr(12),
						Sets:         3,
						Weight:       FloatPtr(150),
						Notes:        "Felt weak today",
						OrderIndex:   1,
					},
					{
						ExerciseName:    "Leg Press",
						Reps:            IntPtr(10),
						DurationSeconds: IntPtr(30),
						Sets:            3,
						Weight:          FloatPtr(200),
						Notes:           "Felt strong today",
						OrderIndex:      2,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := workoutStore.CreateWorkout(tt.workout)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, createdWorkout)
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.DurationMinutes, createdWorkout.DurationMinutes)
			assert.Equal(t, tt.workout.CaloriesBurned, createdWorkout.CaloriesBurned)
			assert.Equal(t, len(tt.workout.Entries), len(createdWorkout.Entries))
			for i, entry := range tt.workout.Entries {
				assert.Equal(t, entry.ExerciseName, createdWorkout.Entries[i].ExerciseName)
				assert.Equal(t, entry.Reps, createdWorkout.Entries[i].Reps)
				assert.Equal(t, entry.Sets, createdWorkout.Entries[i].Sets)
				assert.Equal(t, entry.Weight, createdWorkout.Entries[i].Weight)
				assert.Equal(t, entry.Notes, createdWorkout.Entries[i].Notes)
				assert.Equal(t, entry.OrderIndex, createdWorkout.Entries[i].OrderIndex)
			}

			retrievedWorkout, err := workoutStore.GetWorkoutByID(int64(createdWorkout.ID))
			require.NoError(t, err)
			assert.NotNil(t, retrievedWorkout)
			assert.Equal(t, createdWorkout.ID, retrievedWorkout.ID)
			assert.Equal(t, createdWorkout.Title, retrievedWorkout.Title)
			assert.Equal(t, createdWorkout.Description, retrievedWorkout.Description)
			assert.Equal(t, createdWorkout.DurationMinutes, retrievedWorkout.DurationMinutes)
			assert.Equal(t, createdWorkout.CaloriesBurned, retrievedWorkout.CaloriesBurned)
			assert.Equal(t, len(createdWorkout.Entries), len(retrievedWorkout.Entries))
			for i, entry := range createdWorkout.Entries {
				assert.Equal(t, entry.ExerciseName, retrievedWorkout.Entries[i].ExerciseName)
				assert.Equal(t, entry.Reps, retrievedWorkout.Entries[i].Reps)
				assert.Equal(t, entry.Sets, retrievedWorkout.Entries[i].Sets)
				assert.Equal(t, entry.Weight, retrievedWorkout.Entries[i].Weight)
				assert.Equal(t, entry.Notes, retrievedWorkout.Entries[i].Notes)
				assert.Equal(t, entry.OrderIndex, retrievedWorkout.Entries[i].OrderIndex)
			}
		})
	}
}

func IntPtr(i int) *int {
	return &i
}

func FloatPtr(f float64) *float64 {
	return &f
}
