# golang make file

.DEFAULT_GOAL := build
NAME=ggcache
VERSION=0.0.1

fmt:
	go fmt ./...


vet: fmt ## Run go vet against code
	go vet ./...

test: vet ## Run unit tests
	go test -v ./...

build: fmt
	go build -o ./bin/$(NAME) -v

run: build
	./bin/$(NAME)

follower: build
	./bin/$(NAME) --listenaddr :4000 --leaderaddr :3000

run-client: 
	go run ./client/main.go

.PHONY: build
