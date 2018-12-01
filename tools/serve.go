package main

import (
	"flag"
	"log"
	"net/http"
)

// Config is a collection of the configurable knobs for this app.
type Config struct {
	BindAddr string
	Dir      string
}

func main() {
	config := Config{}

	flag.StringVar(&config.BindAddr, "BindAddr", ":8001", "The bind address for this server.")
	flag.StringVar(&config.Dir, "Dir", "./build", "The directory to serve up.")
	flag.Parse()

	log.Printf("Serving %s on %s", config.Dir, config.BindAddr)
	log.Fatal(http.ListenAndServe(config.BindAddr, http.FileServer(http.Dir(config.Dir))))
}
