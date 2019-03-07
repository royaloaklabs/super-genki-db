# super-genki-db

Quick utility to create the database for the Super Genki: Japanese Dictionary app

## Prerequisites
### go-sqlite3
This utility depends on the [go-sqlite3](https://github.com/mattn/go-sqlite3) driver to write to the SQLite database, it can be installed by running:
```bash
$ go get github.com/mattn/go-sqlite3
```
More information about the project can be found at its [Github page](https://github.com/mattn/go-sqlite3).

### Dictionary and Corpus Files
The main components for creating the database. These needed to be downloaded separately. For more information, view `README` in the `data` folder.

## Building the Database
From within the project root:
```bash
$ go run main.go
```
Wait a bit...

Parsing the XML file and building the database takes a little while.
