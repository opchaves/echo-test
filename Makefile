SHORT_ID := $(shell git rev-parse --short HEAD)
PACKAGES := $(shell go list ./...)
name := $(shell basename ${PWD})
export V_TAG = ${SHORT_ID}

$(eval include .env)
$(eval export)

all: help

.PHONY: help
help: Makefile
	@echo
	@echo " Choose a make command to run"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## run: run the server with "go run"
.PHONY: run
run:
	go run ./cmd/server/main.go

## dev: run the server with "air" and "npm run dev"
.PHONY: dev
dev:
	make watch-api & make dev-web

## dev-web: run the web client with "npm run dev"
.PHONY: dev-web
dev-web:
	cd web && npm run dev

## tidy: run go mod tidy
.PHONY: tidy
tidy:
	go mod tidy

## build: build the server and web client for production
.PHONY: build
build:
	@make build-web
	@make build-api
	@echo "Done!"

## build-api: build the server
.PHONY: build-api
build-api:
	@echo "Building API..."
	@go build -o bin/server ./cmd/server/main.go
	@chmod +x bin/server

## build-web: build the web client
.PHONY: build-web
build-web:
	@echo "Building Web..."
	@cd ./web && npm install && npm run build

## start: start the server binary "bin/server"
.PHONY: start
start:
	APP_ENV=production ./bin/server

## docker-up: start the docker containers
.PHONY: docker-up
docker-up:
	@docker compose up --build -d

## docker-down: store and remove the docker containers
.PHONY: docker-down
docker-down:
	@docker compose down

## docker-build: build the docker image
.PHONY: docker-build
docker-build:
	@docker build --target prod -t ${IMAGE_NAME}:${V_TAG} .

## docker-run: run the docker image
.PHONY: docker-run
docker-run:
	@docker run --rm --name ${APP_NAME} -p ${PORT}:8080 ${IMAGE_NAME}:${V_TAG}

## psql: log into the dev database
.PHONY: psql
psql:
	@docker compose exec postgres psql -U postgres -d ${DB_NAME}

## psql-test: log into the test database
.PHONY: psql-test
psql-test:
	@docker compose exec postgres psql -U postgres -d ${DB_NAME}_test

## db-dump: dump the database schema into db/schema.sql
.PHONY: db-dump
db-dump:
	@docker compose exec postgres pg_dump -U postgres --schema-only -d ${DB_NAME} > db/schema.sql

## db-seed: populate the dev database with initial data
.PHONY: db-seed
db-seed:
	@go run ./cmd/seed/main.go

## migrate: apply all up migrations to dev and test databases
.PHONY: migrate
migrate:
	@make migrate-dev
	@make migrate-test

## migrate-dev: apply all up migrations to dev database
.PHONY: migrate-dev
migrate-dev:
	@migrate -database ${DATABASE_URL} -path ./db/migrations up;\

## migrate-test: apply all up migrations to test database
.PHONY: migrate-test
migrate-test:
	@migrate -database ${DATABASE_URL_TEST} -path ./db/migrations up;\

## migrate-new: generate a new migrate given it name (make migrate-new name=create_some_table)
.PHONY: migrate-new
migrate-new:
	migrate create -ext sql -dir db/migrations -seq $(name);\

## migrate-down: apply all down migrations to dev and test databases
.PHONY: migrate-down
migrate-down:
	@migrate -database ${DATABASE_URL} -path ./db/migrations down
	@migrate -database ${DATABASE_URL_TEST} -path ./db/migrations down

## db-drop: drop the dev and test databases
.PHONY: db-drop
db-drop:
	@migrate -database ${DATABASE_URL} -path ./db/migrations drop
	@migrate -database ${DATABASE_URL_TEST} -path ./db/migrations drop

## db-drop-test: drop the test database
.PHONY: db-drop-test
db-drop-test:
	echo "Droping test database"
	@migrate -database ${DATABASE_URL_TEST} -path ./db/migrations drop

## install-tools: install the required tools (migrate, air, sqlc)
.PHONY: install-tools
install-tools:
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/cosmtrek/air@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

## watch-api: run the server with "air"
.PHONY: watch-api
watch-api:
	@air

## sqlc: generate go files with queries and models
.PHONY: sqlc
sqlc:
	@sqlc generate

## test: run the tests
.PHONY: test
test:
	APP_ENV=test go test ./... -count=1
