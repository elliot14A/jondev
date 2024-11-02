DB_PATH=.jondev/sqlite/jondev.db
MIGRATIONS_PATH=./migrations
BUILD_DIR=build
BINARY_NAME=jondev-server

.PHONY: all clean build run migrate-* sqlc install dev db-init

all: clean build

clean:
	rm -rf ${BUILD_DIR}
	mkdir -p ${BUILD_DIR}

build: build-server build-web

build-server:
	go build -o ${BUILD_DIR}/${BINARY_NAME} main.go

build-web:
	cd web && npm run build

dev-server:
	go run ./cmd/server

dev-web:
	cd web && npm run dev

install:
	cd web && npm install
	go mod download

db-init:
	touch ${DB_PATH}

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir ${MIGRATIONS_PATH} -seq $$name

migrate-up:
	migrate -database "sqlite3://${DB_PATH}" -path ${MIGRATIONS_PATH} up

migrate-down:
	migrate -database "sqlite3://${DB_PATH}" -path ${MIGRATIONS_PATH} down

migrate-force:
	@read -p "Enter version: " version; \
	migrate -database "sqlite3://${DB_PATH}" -path ${MIGRATIONS_PATH} force $$version

migrate-version:
	migrate -database "sqlite3://${DB_PATH}" -path ${MIGRATIONS_PATH} version

sqlc-gen:
	sqlc generate

run: build
	./${BUILD_DIR}/${BINARY_NAME}
