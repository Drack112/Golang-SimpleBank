.PHONY: migrate-create migrate-up migrate-down migrate-force sqlc test postgres server

PWD = $(shell pwd)
PORT = 5432
USER = drack
PASSWORD = 123
NAME = simplebank
HOST = localhost

N = 1

postgres:
	docker-compose up db -d

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
	go test -v -cover ./db/sqlc

server:
	go run main.go
