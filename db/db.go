package db

import (
	"fmt"
	"log"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
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
	fmt.Printf("[DEBUG] Connecting to Database (%s, %s)\n", Driver, Connection)

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

func PopulateDatabase(entries []*SGEntry) (err error) {
	fmt.Println("[INFO] Populating Database")
	// drop and create the virtual table
	_, err = SQL.Exec("DROP TABLE IF EXISTS einihongo")
	if err != nil {
		return
	}

	_, err = SQL.Exec("CREATE VIRTUAL TABLE einihongo USING fts4(japanese,furigana,english,romanji,freq)")
	if err != nil {
		return
	}

	stmt, err := SQL.Prepare("INSERT INTO einihongo(docid,japanese,furigana,english,romanji,freq) VALUES(?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		// execute the statement
		_, err = stmt.Exec(entry.Id, entry.Japanese, entry.Furigana, entry.English, entry.Romanji, entry.Frequency)
		if err != nil {
			return err
		}
	}

	return nil
}
