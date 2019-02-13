# super-genki-db

Quick utility to create the database for the Super Genki: Japanese Dictionary app

## Prerequisites
### go-sqlite3
This utility depends on the [go-sqlite3](https://github.com/mattn/go-sqlite3) driver to write to the SQLite database, it can be installed by running:
```bash
$ go get github.com/mattn/go-sqlite3
```
More information about the project can be found at its [Github page](https://github.com/mattn/go-sqlite3).

### JMDict Japanese Dictionary File
You'll need to download the dictionary file located at the [JMDict Homepage](http://edrdg.org/jmdict/j_jmdict.html). As of right now, this only supports the `JMDict_e` file as we are only supporting English translations. Use the `gzip` utility to extract the XML file and copy it to the `data` directory.

```bash
$ gzip -d JMdict_e.gz
$ cp JMDict_e $GOPATH/src/github.com/Xsixteen/SuperGenki-Utilities/data/.
```

## Building the Database
From within the project root:
```bash
$ go run main.go
```
Wait a bit...

Parsing the XML file and building the database takes a little while.
