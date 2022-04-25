package handlers

import (
	"fmt"
	"os"
)

var POSTGRES_DATABASE string = os.Getenv("VENTRY_POSTGRES_DATABASE")
var POSTGRES_USERNAME string = os.Getenv("VENTRY_POSTGRES_USERNAME")
var POSTGRES_PASSWORD string = os.Getenv("VENTRY_POSTGRES_PASSWORD")
var POSTGRES_HOST string = os.Getenv("VENTRY_POSTGRES_HOST")
var POSTGRES_PORT string = os.Getenv("VENTRY_POSTGRES_PORT")
var POSTGRES_QUERIES string = PostgresQueries()

var POSTGRES_CONNECTION_STRING string = fmt.Sprintf(
	"postgres://%s:%s@%s:%s/%s%s",
	POSTGRES_USERNAME,
	POSTGRES_PASSWORD,
	POSTGRES_HOST,
	POSTGRES_PORT,
	POSTGRES_DATABASE,
	POSTGRES_QUERIES,
)

func PostgresQueries() string {
	queryString := "?"
	if VENTRY_ENV != "PROD" {
		queryString += "sslmode=disable"
	}

	if queryString == "?" {
		queryString = ""
	}

	return queryString
}

type Item struct {
	ItemId   int64
	Title    string
	Quantity int64
	Price    float64
	Owner    string
	Supplier string
	Shipper  string
	Created  string
	Modified string
}
