package main

import (
	"log"

	"github.com/varnit-ta/PlacementLog/cmd/server"
)

func main() {
	_, err := server.InitApp()

	if err != nil {
		log.Fatal(err)
	}
}
