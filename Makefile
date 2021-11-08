ENV_FILE := .env
ENV := $(shell cat $(ENV_FILE))

.PHONY:run
run:
	$(ENV) go run main.go

.PHONY:fmt
fmt:
	go fmt ./...
