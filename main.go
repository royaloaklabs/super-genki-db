package main

import (
	"fmt"
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

	fmt.Println(len(jmdict.Entries))
}
