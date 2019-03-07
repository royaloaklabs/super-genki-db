package main

import (
	"log"

	"github.com/Xsixteen/super-genki-db/db"
	"github.com/Xsixteen/super-genki-db/freq"
	"github.com/Xsixteen/super-genki-db/jmdict"
)

func main() {
	freq.BuildFrequencyData()

	err := jmdict.Parse()
	if err != nil {
		log.Fatal(err)
	}

	databaseEntries := make([]*db.SGEntry, 0)
	for _, entry := range jmdict.Entries {
		databaseEntries = append(databaseEntries, db.NewSGEntryFromJMDict(entry))
	}

	db.Connect()
	err = db.PopulateDatabase(databaseEntries)
	if err != nil {
		log.Fatal(err)
	}
}
