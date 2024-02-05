#!make

ifeq ($(shell test -f ./.env && echo yes), yes)
    include ./.env
    export $(shell sed 's/=.*//' ./.env)
endif

DOCKER_BASE_IMAGE_BUILD := docker build -t deewye/base-image:1.0 ./build/base-image
DOCKER_RUN := docker run --rm -v $(PWD):/go/src -w /go/src deewye/base-image:1.0

## Generating server
docker-base-image-build:
	$(DOCKER_BASE_IMAGE_BUILD)

gen-sql:
	@echo "Removing generated files..."
	@rm -rf ./gen/db/*
	@mkdir -p ./gen/db
	@echo "Generate sql..."
	@sqlc generate

docker-gen-sql:
	$(DOCKER_RUN) make gen-sql

gen-grpc:
	@echo "Removing generated files..."
	@rm -rf ./gen/proto/*
	@rm -rf ./gen/docs/*
	protoc --go_out=./gen \
           --go-grpc_out=./gen \
		   --grpc-gateway_out=./gen \
	       --openapiv2_out ./gen/docs \
	       --openapiv2_opt logtostderr=true,allow_merge=true,merge_file_name=docs \
	       --go_opt=paths=source_relative \
	       --go-grpc_opt=paths=source_relative \
	       --grpc-gateway_opt=paths=source_relative \
		   ./proto/*.proto

docker-gen-grpc:
	$(DOCKER_RUN) make gen-grpc

## Migrations
migrations-create:
	@read -p "Name of the migration: " migration \
	&& echo "Create migrations $$migration at postgres ${USERS_POSTGRES_MASTER_DSN}" \
	&& goose -dir migrations postgres "${USERS_POSTGRES_MASTER_DSN}" create $$migration sql

migrations-up:
	@goose -dir migrations postgres "${USERS_POSTGRES_MASTER_DSN}" up

migrations-down:
	@goose -dir migrations postgres "${USERS_POSTGRES_MASTER_DSN}" down

## Init
install-tools:
	@go install github.com/pressly/goose/v3/cmd/goose@latest

init-dependency:
	@go mod tidy
	@go mod vendor

## Run
run:
	@echo "Running app..."
	@go run -race ./cmd/main.go
