package httpserver

import (
	"gopher/internal/config"
	"net/http"
)

type HttpServer struct {
	config *config.Config
}

func New(config *config.Config) *HttpServer {
	return &HttpServer{
		config: config,
	}
}

func (h HttpServer) Start() error {
	if err := http.ListenAndServe(h.config.Address, nil); err != nil {
		return err
	}

	return nil
}
