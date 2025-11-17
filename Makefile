.PHONY: build

build:
	go build -v ./cmd/http-server

start:
	CONFIG_PATH=./config/local.yaml go run ./cmd/http-server/main.go

.DEFAULT_GOAL := start