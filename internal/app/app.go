package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gbuenodev/goProject/internal/api"
)

type App struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
}

func NewApp() (*App, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	workoutHandler := api.NewWorkoutHandler()

	app := &App{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
	}

	return app, nil
}

func (a *App) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
