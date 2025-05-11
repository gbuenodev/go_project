BINARY_NAME=workout_backend
BUILD_OUTPUT_DIR=./bin
SRC=./

.PHONY: up down restart logs cleanup run

all: build

build: clean
	mkdir -p $(BUILD_OUTPUT_DIR)
	go build -o $(BUILD_OUTPUT_DIR)/$(BINARY_NAME) $(SRC)

run: build
	$(BUILD_OUTPUT_DIR)/$(BINARY_NAME)

clean:
	rm -rf $(BUILD_OUTPUT_DIR)

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

docker-restart: docker-down docker-up

docker-logs:
	docker compose logs -f

docker-clean:
	docker system prune -f --volumes
