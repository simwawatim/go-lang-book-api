-include .env
export

DOCKERHUB_USER ?= simwawa
IMAGE ?= $(DOCKERHUB_USER)/book-api
TAG ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo dev)

.PHONY: tidy run build test migrate-up migrate-down db docker-build docker-push release

tidy:
	go mod tidy

run:
	go run ./cmd/server

build:
	CGO_ENABLED=0 go build -o bin/server ./cmd/server
	CGO_ENABLED=0 go build -o bin/migrate ./cmd/migrate

test:
	go test ./...

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

db:
	docker compose up -d db

docker-build:
	docker build -t $(IMAGE):$(TAG) -t $(IMAGE):latest .

docker-push: docker-build
	docker push $(IMAGE):$(TAG)
	docker push $(IMAGE):latest

release: docker-push
	cd k8s && kustomize edit set image book-api=$(IMAGE):$(TAG)
