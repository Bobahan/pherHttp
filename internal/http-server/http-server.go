package httpserver

import (
	"gopher/internal/config"
	"net/http"

	"github.com/sirupsen/logrus"
)

type HttpServer struct {
	config *config.Config
	logger *logrus.Logger
}

func New(config *config.Config) *HttpServer {
	return &HttpServer{
		config: config,
		logger: logrus.New(),
	}
}

func (h *HttpServer) Start() error {
	if err := h.configureLogger(); err != nil {
		return err
	}

	h.logger.Info("starting http server")

	if err := http.ListenAndServe(h.config.Address, nil); err != nil {
		return err
	}

	return nil
}

func (h *HttpServer) configureLogger() error {
	level, err := logrus.ParseLevel(h.config.LogLevel)

	if err != nil {
		return err
	}

	h.logger.SetLevel(level)

	return nil
}
