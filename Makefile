include app.env

.PHONY: migrate-create migrate-up migrate-down migrate-force sqlc test server postgres-db postgres-db-test

PWD = $(shell pwd)
PORT = ${DB_PORT}
USER = ${DB_USER}
PASSWORD = ${DB_PASSWORD}
NAME = ${DB_NAME}
HOST = ${DB_HOST_TEST}

N = 1

migrate-create:
	@echo "---Creating migration files---"
	migrate create -ext sql -dir $(PWD)/db/migrations -seq -digits 5 $(NAME)

migrate-up:
	migrate -source file://$(PWD)/db/migrations -database postgres://${USER}:${PASSWORD}@${HOST}:$(PORT)/$(NAME)?sslmode=disable up $(N)

migrate-down:
	migrate -source file://$(PWD)/db/migrations -database postgres://${USER}:${PASSWORD}@${HOST}:$(PORT)/$(NAME)?sslmode=disable down $(N)

migrate-force:
	migrate -source file://$(PWD)/db/migrations -database postgres://${USER}:${PASSWORD}@${HOST}:$(PORT)/$(NAME)?sslmode=disable force

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	set -a && source ./app.env && docker-compose up --build -V

postgres-db:
	set -a && source ./app.env && docker-compose up db --build -V

postgres-db-test:
	set -a && source ./app.env && docker-compose up db_test --build -V
