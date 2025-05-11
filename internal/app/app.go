package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gbuenodev/goProject/internal/api"
	"github.com/gbuenodev/goProject/internal/store"
	"github.com/gbuenodev/goProject/migrations"
)

type App struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	DB             *sql.DB
}

func NewApp() (*App, error) {
	dbConfig := store.DBConfig{
		Provider: "Postgres",
		Driver:   "pgx",
		User:     "postgres",
		Password: "postgres",
		DBName:   "postgres",
		Host:     "localhost",
		Port:     5432,
		SSL:      "disable",
	}

	DB, err := store.Open(&dbConfig)
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(DB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	workoutStore := store.NewDBWorkoutStore(DB)

	workoutHandler := api.NewWorkoutHandler(workoutStore)

	app := &App{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		DB:             DB,
	}

	return app, nil
}

func (a *App) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
