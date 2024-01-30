#!make

ifeq ($(shell test -f ./.env && echo yes), yes)
    include ./.env
    export $(shell sed 's/=.*//' ./.env)
endif

PROJECT_NAME:=$(shell basename ${PWD})
BUILD_TIME:=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

export APP_NAME:=$(PROJECT_NAME)

DOCKER_BASE_IMAGE_BUILD := docker build -t deewye/base-image:1.0 ./build/base-image
DOCKER_RUN := docker run --rm -v $(PWD):/go/src -w /go/src deewye/base-image:1.0

## Init
docker-base-image-build:
	$(DOCKER_BASE_IMAGE_BUILD)

# Generating server
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

# Migrations
migrations-create:
	@read -p "Name of the migration: " migration \
	&& echo "Create migrations $$migration at postgres ${USERS_POSTGRES_MASTER_DSN}" \
	&& goose -dir migrations postgres "${USERS_POSTGRES_MASTER_DSN}" create $$migration sql

migrations-up:
	@goose -dir migrations postgres "${USERS_POSTGRES_MASTER_DSN}" up

migrations-down:
	@goose -dir migrations postgres "${USERS_POSTGRES_MASTER_DSN}" down