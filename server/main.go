package main

import (
	"log"
	"os"

	"github.com/tntmeijs/pitona/server/api"
	"github.com/tntmeijs/pitona/server/obdii"
	"golang.org/x/exp/slices"
)

func main() {
	args := os.Args[1:]
	isDebug := false
	port := 80

	if slices.Contains(args, "-d") || slices.Contains(args, "--debug") {
		log.Println("Debugging mode has been enabled")
		isDebug = true
		port = 8080
	}

	instance := obdii.Instance{}
	instance.Connect(isDebug)
	defer instance.Disconnect()

	server := api.NewApiServer(&instance)
	server.Start(isDebug, port)
	log.Println("Application exited successfully")
}
