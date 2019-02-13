package main

import (
	"log"

	"github.com/funayman/super-genki-db/db"
	"github.com/funayman/super-genki-db/jmdict"
)

func main() {
	err := jmdict.Parse()
	if err != nil {
		log.Fatal(err)
	}

	db.Connect()

	err = db.InsertData()
	if err != nil {
		log.Fatal(err)
	}
}
