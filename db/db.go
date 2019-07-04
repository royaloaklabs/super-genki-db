package db

import (
	"fmt"
	"log"
	"strings"

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

	SQL.Exec("DROP TABLE IF EXISTS definitions")
	SQL.Exec("DROP TABLE IF EXISTS readings")
	SQL.Exec("DROP TABLE IF EXISTS sense_misc")

	_, err = SQL.Exec("CREATE VIRTUAL TABLE einihongo USING fts4(japanese,furigana,english,romaji,freq)")
	if err != nil {
		return
	}

	SQL.Exec("CREATE TABLE definitions(id INTEGER PRIMARY KEY AUTOINCREMENT, docid INTEGER, pos TEXT, gloss TEXT)")
	SQL.Exec("CREATE TABLE readings(id INTEGER PRIMARY KEY, japanese TEXT, furigana TEXT, altkanji TEXT, altkana TEXT, romaji TEXT)")
	SQL.Exec("CREATE TABLE sense_misc(senseid INTEGER, docid INTEGER, misc TEXT, PRIMARY KEY (senseid, docid, misc))")

	ftsStmt, err := SQL.Prepare("INSERT INTO einihongo(docid,japanese,furigana,english,romaji,freq) VALUES(?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	definitionStmt, err := SQL.Prepare("INSERT INTO definitions(docid,pos,gloss) VALUES(?,?,?)")
	if err != nil {
		return err
	}

	readingStmt, err := SQL.Prepare("INSERT INTO readings(id,japanese,furigana,altkanji,altkana,romaji) VALUES(?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	miscStmt, err := SQL.Prepare("INSERT INTO sense_misc(senseid, docid, misc) VALUES(?,?,?)")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		// insert into FTS4 entries
		_, err = ftsStmt.Exec(entry.Id, entry.Japanese, entry.Furigana, entry.English, entry.Romaji, entry.Frequency)
		if err != nil {
			return err
		}

		// add all senses into definitions table
		for _, sense := range entry.Sense {
			rslt, err := definitionStmt.Exec(entry.Id, sense.POS, sense.Gloss)
			if err != nil {
				return err
			}
			rowId, err := rslt.LastInsertId()
			if err != nil {
				return err
			}

			if sense.Misc != "" {
				_, err := miscStmt.Exec(rowId, entry.Id, sense.Misc)
				if err != nil {
					return err
				}
			}
		}

		// add all readings
		_, err = readingStmt.Exec(entry.Id, entry.Japanese, entry.Furigana,
			strings.Join(entry.KanjiAlt, " "),
			strings.Join(entry.ReadingAlt, " "),
			entry.Romaji)
		if err != nil {
			return err
		}

	}

	SQL.Exec("CREATE INDEX idx_definitions_docid ON definitions(docid)")
	return nil
}
