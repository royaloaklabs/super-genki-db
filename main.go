package main

import (
	"log"

	"github.com/Xsixteen/super-genki-db/db"
	"github.com/Xsixteen/super-genki-db/jmdict"
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
