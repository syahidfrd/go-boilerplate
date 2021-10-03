# GO Boilerplate

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