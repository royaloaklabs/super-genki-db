package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/funayman/super-genki-db/sql"
)

func main() {
	//get the file
	data, err := os.Open("./data/JMdict_e")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	//parse data
	words, err := LoadJMDict(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(words)
}

func LoadJMDict(f io.Reader) (words []*sql.Entry, err error) {
	d, _ := ioutil.ReadAll(f)

	//needed to fix issue
	//https://groups.google.com/forum/#!topic/golang-nuts/yF9RM9rnkYc
	//get all <!ENTITY> objects in XML
	//fix errors when trying to parse &n; &hon; etc
	var rEntity = regexp.MustCompile(`<!ENTITY\s+([^\s]+)\s+"([^"]+)">`)
	entities := make(map[string]string)
	entityDecoder := xml.NewDecoder(bytes.NewReader(d))
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

	decoder := xml.NewDecoder(bytes.NewReader(d)) //go through the data again
	decoder.Entity = entities                     //load entities into the decoder EntityMap
	for {
		//grab all <entry> tokens and Unmarshal into struct
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "entry" {
				var e *sql.Entry

				if err = decoder.DecodeElement(&e, &se); err != nil {
					return nil, err
				}
				words = append(words, e)
			}
		default:
			//do nothing
		}
	}
	return words, nil
}
