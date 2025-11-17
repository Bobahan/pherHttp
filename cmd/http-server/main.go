package main

import (
	"gopher/internal/config"
	httpserver "gopher/internal/http-server"
	"log"
)

func main() {
	config := config.New()

	s := httpserver.New(config)

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
