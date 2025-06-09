.PHONY: all deps constructor build run clean

PROJECT_DIR=$(shell pwd)
BUILD_DIR=$(PROJECT_DIR)/build
SERVER_HOST?=localhost
SERVER_PORT?=8080
SERVER_URL=http://$(SERVER_HOST):$(SERVER_PORT)

LISTER=listener_net
SERVER=server
OUTPUT_LISTER=$(BUILD_DIR)/$(LISTER)
OUTPUT_SERVER=$(BUILD_DIR)/$(SERVER)

all: constructor deps build
	docker-compose up --build

deps:
	go mod download

constructor:
	mkdir -p $(BUILD_DIR)

build:
	go build -o $(OUTPUT_LISTER) cmd/net/listener_net.go

run: build
	$(OUTPUT_LISTER)

clean:
	rm -rf $(BUILD_DIR)
	docker-compose down
