// Quicknote stores and searches tens of thousands of short notes.
//
// Copyright (C) 2017  Andrew Miller <amiller@amilx.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package sqlite

import (
	"database/sql"
	"errors"
	"reflect"
	"sync"

	// go-sqlite3 must be imported for initialization
	"github.com/anmil/quicknote/note"
	_ "github.com/mattn/go-sqlite3"
)

var schema = `PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS books (
	id       INTEGER   PRIMARY KEY AUTOINCREMENT,
	created  TIMESTAMP NOT NULL,
	modified TIMESTAMP NOT NULL,
	name     TEXT UNIQUE
);

CREATE INDEX IF NOT EXISTS index_books_name ON books (name);
CREATE INDEX IF NOT EXISTS index_books_created ON books (created);

CREATE TABLE IF NOT EXISTS notes (
	id       INTEGER   PRIMARY KEY AUTOINCREMENT,
	created  TIMESTAMP NOT NULL,
	modified TIMESTAMP NOT NULL,
	bk_id    INTEGER   NOT NULL REFERENCES books(id) ON DELETE CASCADE,
	type     TEXT      NOT NULL,
	title    TEXT,
	body     TEXT
);

CREATE INDEX IF NOT EXISTS index_notes_bk ON notes (bk_id);
CREATE INDEX IF NOT EXISTS index_notes_type ON notes (type);
CREATE INDEX IF NOT EXISTS index_notes_created ON notes (created);
CREATE INDEX IF NOT EXISTS index_notes_bk_created ON notes (bk_id, created);
CREATE INDEX IF NOT EXISTS index_notes_type_created ON notes (type, created);

CREATE TABLE IF NOT EXISTS tags (
	id       INTEGER   PRIMARY KEY AUTOINCREMENT,
	created  TIMESTAMP NOT NULL,
	modified TIMESTAMP NOT NULL,
	name     TEXT UNIQUE
);

CREATE INDEX IF NOT EXISTS index_tags_name ON tags (name);

CREATE TABLE IF NOT EXISTS note_tag (
	note_id INTEGER REFERENCES notes(id) ON DELETE CASCADE,
	tag_id  INTEGER REFERENCES tags(id) ON DELETE CASCADE,
	PRIMARY KEY (note_id, tag_id)
);

CREATE INDEX IF NOT EXISTS index_note_tag_note_id ON note_tag (note_id);
CREATE INDEX IF NOT EXISTS index_note_tag_tag_id ON note_tag (tag_id);

CREATE TABLE IF NOT EXISTS note_book_tag (
	note_id INTEGER REFERENCES notes(id) ON DELETE CASCADE,
	bk_id   INTEGER REFERENCES books(id) ON DELETE CASCADE,
	tag_id  INTEGER REFERENCES tags(id) ON DELETE CASCADE,
	PRIMARY KEY (note_id, bk_id, tag_id)
);`

// Maximum number of wild-card variables SQlite can parse
const sqliteMaxVariableNumber = 999

// ErrInvalidArguments invalid arguments were given
var ErrInvalidArguments = errors.New("Invalid arguments given to SQLite database")

// Database provides an interface to SQLite
type Database struct {
	db     *sql.DB
	mux    *sync.Mutex
	DBPath string

	tagNameCache  map[string]*note.Tag
	bookNameCache map[string]*note.Book
}

// NewDatabase returns a data Database
func NewDatabase(dbPath ...string) (*Database, error) {
	if len(dbPath) != 1 {
		return nil, ErrInvalidArguments
	}

	db, err := sql.Open("sqlite3", dbPath[0])
	if err != nil {
		return nil, err
	}

	// db.SetMaxIdleConns(1)
	// db.SetMaxOpenConns(1)

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return &Database{
		db:            db,
		mux:           &sync.Mutex{},
		DBPath:        dbPath[0],
		tagNameCache:  make(map[string]*note.Tag),
		bookNameCache: make(map[string]*note.Book),
	}, nil
}

// GetTableNames returns a list of all table names
func (d *Database) GetTableNames() ([]string, error) {
	d.mux.Lock()
	defer d.mux.Unlock()

	sqlStr := "SELECT name FROM sqlite_master WHERE type='table' ORDER BY name;"

	rows, err := d.db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tables := make([]string, 0)
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}

		tables = append(tables, name)
	}

	return tables, nil
}

// Close closes the database
func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) getTxStmt(sqlStmt string) (*sql.Tx, *sql.Stmt, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, nil, err
	}

	stmt, err := tx.Prepare(sqlStmt)
	if err != nil {
		return nil, nil, err
	}

	return tx, stmt, nil
}

// splitSliceToChuck slice s into chucks containing the maximum number of
// objects we can use in a statement.
//
// wcPreRecord is the number of wild cards each object requires
func splitSliceToChuck(s interface{}, wcPreRecord int) ([]interface{}, error) {
	// I'm not fund of using reflect as I believe it hurts performance, but I don't
	// know any other way of doing this and still be usable by multiple types
	slice := reflect.ValueOf(s)
	if slice.Kind() != reflect.Slice {
		return nil, errors.New("arg is not a slice")
	}

	maxChuck := sqliteMaxVariableNumber / wcPreRecord
	chuck := make([]interface{}, 0)

	numRecords := slice.Len()
	for i := 0; i < numRecords; i = i + maxChuck {
		end := i + maxChuck
		if end > numRecords {
			end = numRecords
		}
		chuck = append(chuck, slice.Slice(i, end).Interface())
	}

	return chuck, nil
}
