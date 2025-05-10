package routes

import (
	"github.com/gbuenodev/goProject/internal/app"
	"github.com/go-chi/chi/v5"
)

func Routes(app *app.App) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)
	r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkoutByID)

	r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)

	return r
}
