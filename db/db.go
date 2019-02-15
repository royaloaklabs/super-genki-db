package db

import (
	"log"
	"strings"

	"database/sql"

	"github.com/Xsixteen/super-genki-db/jmdict"
	_ "github.com/mattn/go-sqlite3"
)

const (
	Delimiter      = " "
	SenseDelimiter = "(SG)"
	GlossDelimiter = ";;"
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

func InsertData() (err error) {
	tx, err := SQL.Begin()
	if err != nil {
		return
	}

	// create the virtual table
	_, err = tx.Exec("CREATE VIRTUAL TABLE einihongo USING fts4(kanji,kana,gloss)")
	if err != nil {
		return
	}

	tx.Commit()
	if err != nil {
		return
	}

	for _, entry := range jmdict.Entries {
		var kanji, kana, sense string

		stmt, err := SQL.Prepare("INSERT INTO einihongo(docid,kanji,kana,gloss) VALUES(?,?,?,?)")
		if err != nil {
			return err
		}

		// compress kanji
		var tempKanji []string
		for _, kanji := range entry.KEle {
			tempKanji = append(tempKanji, kanji.Keb)
		}
		kanji = strings.Join(tempKanji, Delimiter)

		// compress kana
		var tempKana []string
		for _, kana := range entry.Rele {
			tempKana = append(tempKana, kana.Reb)
		}
		kana = strings.Join(tempKana, Delimiter)

		// compress sense
		var tempGloss []string
		var tempSense []string
		for _, sense := range entry.Sense {
			for _, gloss := range sense.Gloss {
				tempGloss = append(tempGloss, gloss.Value)
			}
			tempSense = append(tempSense, strings.Join(tempGloss, GlossDelimiter))
		}
		sense = strings.Join(tempSense, SenseDelimiter)

		// execute the statement
		_, err = stmt.Exec(entry.EntSeq, kanji, kana, sense)
		if err != nil {
			return err
		}
	}

	return nil
}
