ifneq (,$(wildcard ./.env))
    include .env
    export
endif

migration-create:
	migrate create -ext sql -dir migration -seq $(name)

migration-up:
	migrate -path migration -database "${DATABASE_URL}" up

migration-down:
	migrate -path migration -database "${DATABASE_URL}" down

run-server:
	go run ./cmd/api/main.go

build-api:
	go build ./cmd/api/main.go

test:
	go test -v ./...

mock:
	mockery --all

.PHONY:
	migration-create
	migration-up
	migration-down
	run-server
	build-api
	test
	mock