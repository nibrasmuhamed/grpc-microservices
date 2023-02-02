package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct {
}

func main() {
	// setting up application configuration.
	// this struct will be used for global configuration for this application.
	app := Config{}

	fmt.Println("Serving on Port : 80")

	// creating a server object and
	// assigning the handlers and port address.
	// as this service runs on docker, it is ok
	// to use port 80.
	srv := &http.Server{
		Addr:    ":80",
		Handler: app.routes(),
	}

	// starting server and checking errors.
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("error while starting up server: %v \n", err)
	}

}
