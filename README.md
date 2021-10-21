# GO Boilerplate
[![Go Reference](https://pkg.go.dev/badge/github.com/andhikayuana/qiscus-unofficial-go.svg)](https://pkg.go.dev/github.com/syahidfrd/go-boilerplate)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/andhikayuana/qiscus-unofficial-go)](https://goreportcard.com/report/github.com/syahidfrd/go-boilerplate)

### Prerequisite
Install go-migrate `https://github.com/golang-migrate/migrate` for running migration.
App requires 2 database (postgreSQL and redis server), run from your local machine or run it using docker with the following command:
```
# run postgreSQL
docker run -d -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=go-boilerplate postgres

# run redis
docker run -d -p 6379:6379 redis
``` 

### Migration
Run below command to run migration
```
migrate -path migration -database "${DATABASE_URL}" up
```

To create a new migration file
```
migrate create -ext sql -dir migration -seq name
```

### Test
Run below command to run test, and make sure that all tests are passing.
```
go test -v ./...
```

### Running
Run below command to run app
```
go run ./cmd/api/main.go
```

You can find usefull commands in `Makefile`.