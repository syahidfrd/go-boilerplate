package datastore

import (
	"database/sql"
	"net/url"
	"time"

	_ "github.com/lib/pq"
)

func NewDatabase(databaseURL string) *sql.DB {
	parseDBUrl, _ := url.Parse(databaseURL)
	db, err := sql.Open(parseDBUrl.Scheme, databaseURL)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(5)

	return db
}
