BINARY_NAME=fizz
VERSION=$(shell git rev-parse --short HEAD)

all: test build
build:
	go build -buildvcs=false -o ${BINARY_NAME} ./cmd/${BINARY_NAME}/
test:
	go test ./...
clean:
	go clean
deps:
	go mod vendor
	go mod tidy
up:
	docker compose -p ${BINARY_NAME} up
