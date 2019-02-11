package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/funayman/super-genki-db/db"
	"github.com/funayman/super-genki-db/jmdict"
)

func main() {
	//get the file
	data, err := os.Open("./data/JMdict_e")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	//parse data
	var words []*jmdict.Entry
	if err != nil {
		log.Fatal(err)
	}

	//needed to fix issue
	//https://groups.google.com/forum/#!topic/golang-nuts/yF9RM9rnkYc
	//get all <!ENTITY> objects in XML
	//fix errors when trying to parse &n; &hon; etc
	var rEntity = regexp.MustCompile(`<!ENTITY\s+([^\s]+)\s+"([^"]+)">`)
	entities := make(map[string]string)
	entityDecoder := xml.NewDecoder(data)
	for {
		t, _ := entityDecoder.Token()
		if t == nil {
			break
		}

		dir, ok := t.(xml.Directive)
		if !ok {
			continue
		}

		for _, m := range rEntity.FindAllSubmatch(dir, -1) {
			entities[string(m[1])] = string(m[2])
		}
	}

	data.Seek(0, 0)
	decoder := xml.NewDecoder(data) //go through the data again
	decoder.Entity = entities       //load entities into the decoder EntityMap
	for {
		//grab all <entry> tokens and Unmarshal into struct
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "entry" {
				var e *jmdict.Entry

				if err = decoder.DecodeElement(&e, &se); err != nil {
					log.Fatal(err)
				}
				words = append(words, e)
			}
		default:
			//do nothing
		}
	}

	db.Connect()

	for _, word := range words {
		fmt.Println(word)
	}
}
