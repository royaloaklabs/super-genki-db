package db

import (
	"log"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	Delimeter = ";"
)

var (
	//SQL is a wrapper for database/sql
	SQL *sql.DB

	//Driver is the database type
	Driver = "sqlite3"

	//Connection to the database
	Connection = "./jisho-main.db"
)

//Connect to database of choice
func Connect() {
	var err error
	SQL, err = sql.Open(Driver, Connection)
	if err != nil {
		log.Fatal("SQL Open error: ", err)
	}

	//we good?
	if err = SQL.Ping(); err != nil {
		log.Fatal("Database connection error: ", err)
	}
}
