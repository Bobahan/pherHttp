package main

import (
	"gopher/internal/config"
	httpserver "gopher/internal/http-server"
	"log"
)

func main() {
	config := config.New()

	if err := httpserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
