package database

import (
	"database/sql"
	"log"
)

type database struct {
	db *sql.DB
}

func New() *database {
	return &database{}
}

func (d *database) GetDB() *sql.DB {
	var err error
	if d.db, err = sql.Open("ramsql", "goimdb"); err != nil {
		log.Fatal(err)
	}
	if err = d.db.Ping(); err != nil {
		log.Fatal(err)
	}

	return d.db
}
