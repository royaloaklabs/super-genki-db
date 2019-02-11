package db

import (
	"log"

	"database/sql"
)

var (
	//SQL is a wrapper for database/sql
	SQL *sql.DB

	//Driver is the database type
	Driver = "sqlite3"

	//Connection to the database
	Connection = "./jisho-main.db"
)

//Setup the SQL connection
type Setup struct {
	Driver     string `json:"driver"`
	Connection string `json:"connect"`
}

//Connect to database of choice
func Connect(info Setup) {
	var err error
	SQL, err = sql.Open(info.Driver, info.Connection)
	if err != nil {
		log.Fatal("SQL Open error: ", err)
	}

	//we good?
	if err = SQL.Ping(); err != nil {
		log.Fatal("Database connection error: ", err)
	}
}
