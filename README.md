# super-genki-db

Quick utility to create the database for the Super Genki: Japanese Dictionary app

# Prerequisites
## go-sqlite3
This utility depends on the [go-sqlite3](https://github.com/mattn/go-sqlite3) driver to write to the SQLite database, it can be installed by running:
```bash
$ go get github.com/mattn/go-sqlite3
```
More information about the project can be found at its [Github page](https://github.com/mattn/go-sqlite3).

## Dictionary and Corpus Files
The main components for creating the database. These needed to be downloaded separately. For more information, view `README` in the `data` folder.

# Building the Database
From within the project root:
```bash
$ go run main.go
```

# Schema
```
System Information
========================================================================

generated by                              SchemaCrawler 15.06.01
generated on                              2019-07-08 12:49:56



Tables
========================================================================



definitions                                                      [table]
------------------------------------------------------------------------
  id                                INTEGER
                                    auto-incremented
  entryid                           INTEGER
  pos                               TEXT
  gloss                             TEXT

Primary Key

                                                           [primary key]
  id                                ascending

Foreign Keys

                                           [foreign key, with no action]
  id <--(0..many) sense_misc.senseid

                                           [foreign key, with no action]
  entryid (0..many)--> einihongo.entryid

                                           [foreign key, with no action]
  pos (0..many)--> entity_members.abbvr

Indexes

idx_definitions_entryid                               [non-unique index]
  entryid                           unknown



dirty_talk                                                        [view]
------------------------------------------------------------------------
  entryid                           INTEGER



einihongo                                                       [vTable]
------------------------------------------------------------------------
  entryid
  japanese
  furigana
  english
  romaji
  freq

Foreign Keys

                                           [foreign key, with no action]
  entryid <--(0..many) definitions.entryid

                                           [foreign key, with no action]
  entryid <--(0..1) readings.entryid

                                           [foreign key, with no action]
  entryid <--(0..many) sense_misc.entryid



entity_members                                                   [table]
------------------------------------------------------------------------
  abbvr                             TEXT
  meaning                           TEXT

Primary Key

                                                           [primary key]
  abbvr                             ascending

Foreign Keys

                                           [foreign key, with no action]
  abbvr <--(0..many) definitions.pos

                                           [foreign key, with no action]
  abbvr <--(0..many) sense_misc.misc

Indexes

sqlite_autoindex_entity_members_1                         [unique index]
  abbvr                             unknown



readings                                                         [table]
------------------------------------------------------------------------
  entryid                           INTEGER
  japanese                          TEXT
  furigana                          TEXT
  altkanji                          TEXT
  altkana                           TEXT
  romaji                            TEXT

Primary Key

                                                           [primary key]
  entryid                           ascending

Foreign Keys

                                           [foreign key, with no action]
  entryid (0..1)--> einihongo.entryid



sense_misc                                                       [table]
------------------------------------------------------------------------
  senseid                           INTEGER
  entryid                           INTEGER
  misc                              TEXT

Primary Key

                                                           [primary key]
  senseid                           ascending
  entryid                           ascending
  misc                              ascending

Foreign Keys

                                           [foreign key, with no action]
  senseid (0..many)--> definitions.id

                                           [foreign key, with no action]
  entryid (0..many)--> einihongo.entryid

                                           [foreign key, with no action]
  misc (0..many)--> entity_members.abbvr

Indexes

sqlite_autoindex_sense_misc_1                             [unique index]
  senseid                           unknown
  entryid                           unknown
  misc                              unknown
```
