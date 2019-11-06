package main

import (
	"flag"
	"github.com/google/gops/agent"
	"github.com/krazybee/cmd/searchsvc"
	"log"
)


func main() {
	opts := agent.Options{
		ShutdownCleanup: true,
	}
	if err := agent.Listen(opts); err != nil {
		log.Fatal(err)
	}

	enableSearchService := flag.Bool("search-service", false, "to run search service")
	configPath := flag.String("config", "/etc/krazybee", "config flag provides the path to the configuration file which should be used with the server")
	flag.Parse()

	if *enableSearchService{
		searchsvc.Initialize(configPath)
	} else {
		log.Println("[Main] use flags: --search-service ")
	}
	log.Println("[Main] exiting..")
}

