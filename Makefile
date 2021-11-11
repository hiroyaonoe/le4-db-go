ENV_FILE := .env
ENV := $(shell cat $(ENV_FILE))

.PHONY:run
run:
	$(ENV) go run main.go

.PHONY:fmt
fmt:
	go fmt ./...

.PHONY:init-db
init-db:
	cat sql/create.sql sql/view.sql sql/insert.sql | PGPASSWORD=postgres psql -h localhost -U postgres le4db

