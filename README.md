# GO Boilerplate
[![Go Reference](https://pkg.go.dev/badge/github.com/andhikayuana/qiscus-unofficial-go.svg)](https://pkg.go.dev/github.com/syahidfrd/go-boilerplate)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/andhikayuana/qiscus-unofficial-go)](https://goreportcard.com/report/github.com/syahidfrd/go-boilerplate)

### Prerequisite

Install go-migrate `https://github.com/golang-migrate/migrate` for running migration

### Migration

Run below command to run migration

```
migrate -path migration -database "${DATABASE_URL}" up
```

To create a new migration file

```
migrate create -ext sql -dir migration -seq name
```

or execute command in `Makefile`

### Running

Run below command to run app

```
go run ./cmd/api/main.go
```