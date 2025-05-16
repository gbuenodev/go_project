package routes

import (
	"github.com/gbuenodev/goProject/internal/app"
	"github.com/go-chi/chi/v5"
)

func Routes(app *app.App) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.Auth)

		// AUTHENTICATED ROUTES
		// WORKOUT ROUTES
		r.Get("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleGetWorkoutByID))
		r.Post("/workouts", app.Middleware.RequireUser(app.WorkoutHandler.HandleCreateWorkout))
		r.Put("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleUpdateWorkoutByID))
		r.Delete("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleDeleteWorkoutByID))

	})

	// HEALTH CHECK
	r.Get("/health", app.HealthCheck)

	// USER ROUTES
	r.Post("/users/register", app.UserHandler.HandleRegisterUser)
	r.Post("/auth", app.TokenHandler.HandleCreateToken)

	return r
}
