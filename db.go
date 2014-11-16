package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func NewDB(url string) (*sqlx.DB, error) {
	fmt.Printf("%v", url)
	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(2)

	return db, err
}
