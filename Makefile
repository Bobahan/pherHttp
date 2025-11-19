.PHONY: build

build:
	CONFIG_PATH=./config/local.yaml go build -v ./cmd/http-server

.PHONY: start
start:
	CONFIG_PATH=./config/local.yaml go run ./cmd/http-server/main.go

.PHONY: test
test: 
	CONFIG_PATH=./config/local.yaml go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := start