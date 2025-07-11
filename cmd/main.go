package main

import (
	"log"
	"net/http"

	"github.com/varnit-ta/PlacementLog/cmd/server"
)

const port = ":8080"

func main() {
	app, err := server.InitApp()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting Server...")

	if err = http.ListenAndServe(port, app.Routes()); err != nil {
		log.Fatalf("error starting the server: %v\n", err)
	}
}
