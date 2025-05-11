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
	flag.IntVar(&port, "port", 8080, "GO backend server port")
	flag.Parse()

	app, err := app.NewApp()
	if err != nil {
		panic(err)
	}

	defer app.DBConn.Close()

	app.Logger.Printf("App is running on port %d\n", port)

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
		app.Logger.Fatal(err)
	}
}
