.PHONY: all deps

PROJECT_DIR=$(shell pwd)
BUILD_DIR=$(PROJECT_DIR)/build
BIN_DIR=$(PROJECT_DIR)/bin
PROGRAMM=$(PROJECT_DIR)/cmd/net/listener_net.go $(PROJECT_DIR)/cmd/server/server.go
SERVER_HOST?=localhost
SERVER_PORT?=8080
SERVER_URL=http://$(SERVER_HOST):$(SERVER_PORT)

APP_NAME=listener_net
OUTPUT=$(BIN_DIR)/$(APP_NAME)

all: constructor deps build
deps:
	go mod download
	go get -u github.com/jteeuwen/go-bindata/...

constructor:
	mkdir -p $(BUILD_DIR) $(BIN_DIR)

build: $(PROGRAMM)
	go build -o $(OUTPUT) $(PROGRAMM)

run: build
	$(OUTPUT)

clean:
	rm -rf $(BUILD_DIR) $(BIN_DIR)

