include .env
.EXPORT_ALL_VARIABLES:

.PHONY:run
run:
	go run main.go

.PHONY:fmt
fmt:
	go fmt ./...

.PHONY:init-db
init-db:
	cat sql/create.sql sql/view.sql sql/insert.sql | PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -U $(DB_USER) $(DB_NAME)

