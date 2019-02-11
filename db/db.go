package db

import (
	"log"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var (
	//SQL is a wrapper for database/sql
	DB *sql.DB

	//Driver is the database type
	Driver = "sqlite3"

	//Connection to the database
	Connection = "./jisho-main.db"
)

//Connect to database of choice
func Connect() {
	var err error
	DB, err = sql.Open(Driver, Connection)
	if err != nil {
		log.Fatal("SQL Open error: ", err)
	}

	//we good?
	if err = DB.Ping(); err != nil {
		log.Fatal("Database connection error: ", err)
	}
}
