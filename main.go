package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/gbuenodev/goProject/internal/app"
	"github.com/gbuenodev/goProject/internal/routes"
)

func main() {
	var port int
	var logLevel string

	flag.IntVar(&port, "port", 8080, "GO backend server port")
	flag.StringVar(&logLevel, "level", "info", "Log Level for the app")
	flag.Parse()

	app, err := app.NewApp(logLevel)
	if err != nil {
		panic(err)
	}

	defer app.DBConn.Close()

	fmt.Printf(`
App is running on port: %d
Log level: %s

`, port, logLevel)

	r := routes.Routes(app)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
