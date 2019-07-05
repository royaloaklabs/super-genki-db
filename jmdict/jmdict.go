package jmdict

import (
	"fmt"
	"os"
	"regexp"

	"encoding/xml"
)

var Entries []*Entry
var XmlEntities map[string]string
var XmlReverseEntities map[string]string // make use of later

func Parse() error {
	fmt.Println("[INFO] Parsing JMDict Data")

	//get the file
	data, err := os.Open("./data/JMdict_e")
	if err != nil {
		return err
	}
	defer data.Close()

	//needed to fix issue
	//https://groups.google.com/forum/#!topic/golang-nuts/yF9RM9rnkYc
	//get all <!ENTITY> objects in XML
	//fix errors when trying to parse &n; &hon; etc
	var rEntity = regexp.MustCompile(`<!ENTITY\s+([^\s]+)\s+"([^"]+)">`)
	XmlEntities = make(map[string]string)
	XmlReverseEntities = make(map[string]string)
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
			XmlEntities[string(m[1])] = string(m[2])
			XmlReverseEntities[string(m[2])] = string(m[1])
		}
	}

	data.Seek(0, 0)
	decoder := xml.NewDecoder(data) //go through the data again
	decoder.Entity = XmlEntities    //load entities into the decoder EntityMap
	for {
		//grab all <entry> tokens and Unmarshal into struct
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "entry" {
				var e *Entry

				if err = decoder.DecodeElement(&e, &se); err != nil {
					return err
				}

				for _, sense := range e.Sense {
					for i, pos := range sense.Pos {
						sense.Pos[i] = XmlReverseEntities[pos]
					}

					for i, misc := range sense.Misc {
						sense.Misc[i] = XmlReverseEntities[misc]
					}
				}

				Entries = append(Entries, e)
			}
		default:
			//do nothing
		}
	}

	return nil
}
