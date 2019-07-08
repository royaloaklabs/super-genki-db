package db

import (
	"fmt"
	"log"
	"strings"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/royaloaklabs/super-genki-db/jmdict"
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
	fmt.Println("[DEBUG] DROP tables")
	if _, err = SQL.Exec("DROP TABLE IF EXISTS einihongo"); err != nil {
		return
	}
	if _, err := SQL.Exec("DROP TABLE IF EXISTS entity_members"); err != nil {
		return err
	}
	if _, err := SQL.Exec("DROP TABLE IF EXISTS definitions"); err != nil {
		return err
	}
	if _, err := SQL.Exec("DROP TABLE IF EXISTS readings"); err != nil {
		return err
	}
	if _, err := SQL.Exec("DROP TABLE IF EXISTS sense_misc"); err != nil {
		return err
	}

	// create all tables
	fmt.Println("[DEBUG] CREATE tables")
	if _, err = SQL.Exec("CREATE VIRTUAL TABLE einihongo USING fts4(entryid,japanese,furigana,english,romaji,freq)"); err != nil {
		return
	}
	if _, err := SQL.Exec("CREATE TABLE entity_members(abbvr TEXT PRIMARY KEY, meaning TEXT)"); err != nil {
		return err
	}
	if _, err := SQL.Exec("CREATE TABLE definitions(id INTEGER PRIMARY KEY AUTOINCREMENT, entryid INTEGER, pos TEXT, gloss TEXT, FOREIGN KEY(entryid) REFERENCES einihongo(entryid), FOREIGN KEY(pos) REFERENCES entity_members(abbvr))"); err != nil {
		return err
	}
	if _, err := SQL.Exec("CREATE TABLE readings(entryid INTEGER PRIMARY KEY, japanese TEXT, furigana TEXT, altkanji TEXT, altkana TEXT, romaji TEXT, FOREIGN KEY(entryid) REFERENCES einihongo(entryid))"); err != nil {
		return err
	}
	if _, err := SQL.Exec("CREATE TABLE sense_misc(senseid INTEGER, entryid INTEGER, misc TEXT, PRIMARY KEY (senseid, entryid, misc), FOREIGN KEY(entryid) REFERENCES einihongo(entryid), FOREIGN KEY(senseid) REFERENCES definitions(id), FOREIGN KEY(misc) REFERENCES entity_members(abbvr))"); err != nil {
		return err
	}

	// begin transaction
	tx, err := SQL.Begin()
	if err != nil {
		return err
	}
	fmt.Println("[DEBUG] INSERT XmlEntities")
	for abbvr, meaning := range jmdict.XmlEntities {
		if _, err := tx.Exec("INSERT INTO entity_members(abbvr, meaning) VALUES(?,?)", abbvr, meaning); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	tx, err = SQL.Begin()
	if err != nil {
		return err
	}
	// prepare statements for tables
	ftsStmt, err := tx.Prepare("INSERT INTO einihongo(entryid,japanese,furigana,english,romaji,freq) VALUES(?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	definitionStmt, err := tx.Prepare("INSERT INTO definitions(entryid,pos,gloss) VALUES(?,?,?)")
	if err != nil {
		return err
	}

	readingStmt, err := tx.Prepare("INSERT INTO readings(entryid,japanese,furigana,altkanji,altkana,romaji) VALUES(?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	miscStmt, err := tx.Prepare("INSERT INTO sense_misc(senseid, entryid, misc) VALUES(?,?,?)")
	if err != nil {
		return err
	}

	fmt.Println("[DEBUG] INSERT entries")
	for _, entry := range entries {
		// insert into FTS4 entries
		_, err = ftsStmt.Exec(entry.Id, entry.Japanese, entry.Furigana, entry.English, entry.Romaji, entry.Frequency)
		if err != nil {
			tx.Rollback()
			return err
		}

		// add all senses into definitions table
		for _, sense := range entry.Sense {
			rslt, err := definitionStmt.Exec(entry.Id, sense.POS, sense.Gloss)
			if err != nil {
				tx.Rollback()
				return err
			}
			rowId, err := rslt.LastInsertId()
			if err != nil {
				tx.Rollback()
				return err
			}

			if sense.Misc.String != "" {
				_, err := miscStmt.Exec(rowId, entry.Id, sense.Misc)
				if err != nil {
					tx.Rollback()
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
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	SQL.Exec("CREATE INDEX idx_definitions_entryid ON definitions(entryid)")

	SQL.Exec("DROP VIEW IF EXISTS dirty_talk")
	SQL.Exec("CREATE VIEW dirty_talk AS SELECT DISTINCT entryid FROM sense_misc WHERE misc = 'vulg' OR misc = 'sl' OR misc = 'm-sl'")
	return nil
}
