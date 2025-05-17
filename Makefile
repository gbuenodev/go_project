BINARY_NAME=workout_server
BUILD_OUTPUT_DIR=./bin
SRC=./
ENV ?= prd
LEVEL ?= INFO
PORT ?= 8080

.PHONY: build run clean test docker-up docker-down docker-restart docker-logs docker-clean stop

build: clean
	mkdir -p $(BUILD_OUTPUT_DIR)
	go build -o $(BUILD_OUTPUT_DIR)/$(BINARY_NAME) $(SRC)

run: docker-up build
	$(BUILD_OUTPUT_DIR)/$(BINARY_NAME) --level=$(LEVEL) --port=$(PORT)

clean:
	rm -rf $(BUILD_OUTPUT_DIR)

test: 
	$(MAKE) docker-up ENV=stg
	go test -v ./...
	$(MAKE) docker-down ENV=stg

docker-up:
	docker compose up --build -d ${ENV}_db

docker-down:
	docker compose stop ${ENV}_db
	docker compose rm -f ${ENV}_db

docker-restart: docker-down docker-up

docker-logs:
	docker compose logs -f ${ENV}_db

docker-clean:
	docker system prune -f --volumes

stop: docker-down clean
