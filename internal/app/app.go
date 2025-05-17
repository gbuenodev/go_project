package app

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/gbuenodev/goProject/internal/api"
	"github.com/gbuenodev/goProject/internal/errors"
	"github.com/gbuenodev/goProject/internal/middleware"
	"github.com/gbuenodev/goProject/internal/store"
	"github.com/gbuenodev/goProject/internal/utils"
	"github.com/gbuenodev/goProject/migrations"
)

type App struct {
	Logger         *slog.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokenHandler
	Middleware     middleware.UserMiddleware
	DBConn         *sql.DB
}

func NewApp(logLevel string) (*App, error) {
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

	DBConn, err := store.Open(&dbConfig)
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(DBConn, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	loggerOpts := errors.SetLoggerLevel(logLevel)
	logger := errors.SetupDefaultLogger(loggerOpts)

	workoutStore := store.NewPostgresWorkoutStore(DBConn)
	userStore := store.NewPostgresUserStore(DBConn)
	tokenStore := store.NewPostgresTokenStore(DBConn)

	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)
	middlewareHandler := middleware.UserMiddleware{UserStore: userStore}

	app := &App{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		UserHandler:    userHandler,
		TokenHandler:   tokenHandler,
		Middleware:     middlewareHandler,
		DBConn:         DBConn,
	}

	return app, nil
}

func (a *App) HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"alive": "true"})
}
