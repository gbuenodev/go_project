package routes

import (
	"github.com/gbuenodev/goProject/internal/app"
	"github.com/go-chi/chi/v5"
)

func Routes(app *app.App) *chi.Mux {
	r := chi.NewRouter()

	// HEALTH CHECK
	r.Get("/health", app.HealthCheck)

	// WORKOUTS ROUTES
	r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
	r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
	r.Put("/workouts/{id}", app.WorkoutHandler.HandleUpdateWorkoutByID)
	r.Delete("/workouts/{id}", app.WorkoutHandler.HandleDeleteWorkoutByID)

	// USERS ROUTES
	r.Post("/users/register", app.UserHandler.HandleRegisterUser)

	return r
}
