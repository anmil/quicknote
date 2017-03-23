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

package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	// pq must be imported for initialization
	_ "github.com/lib/pq"
)

var schema = `
CREATE TABLE IF NOT EXISTS books (
	id       SERIAL   PRIMARY KEY,
	created  TIMESTAMP NOT NULL,
	modified TIMESTAMP NOT NULL,
	name     TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS notes (
	id       SERIAL   PRIMARY KEY,
	created  TIMESTAMP NOT NULL,
	modified TIMESTAMP NOT NULL,
	bk_id    INTEGER   NOT NULL REFERENCES books(id) ON DELETE CASCADE,
	type     TEXT      NOT NULL,
	title    TEXT,
	body     TEXT
);

CREATE INDEX IF NOT EXISTS idx_notes_bk_id ON notes (bk_id);

CREATE TABLE IF NOT EXISTS tags (
	id       SERIAL   PRIMARY KEY,
	created  TIMESTAMP NOT NULL,
	modified TIMESTAMP NOT NULL,
	name     TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS note_tag (
	note_id INTEGER REFERENCES notes(id) ON DELETE CASCADE,
	tag_id  INTEGER REFERENCES tags(id) ON DELETE CASCADE,
	PRIMARY KEY (note_id, tag_id)
);

CREATE TABLE IF NOT EXISTS note_book_tag (
	note_id INTEGER REFERENCES notes(id) ON DELETE CASCADE,
	bk_id   INTEGER REFERENCES books(id) ON DELETE CASCADE,
	tag_id  INTEGER REFERENCES tags(id) ON DELETE CASCADE,
	PRIMARY KEY (note_id, bk_id, tag_id)
);`

// ErrInvalidArguments invalid arguments were given
var ErrInvalidArguments = errors.New("Invalid arguments given to PostgreSQL database")

var maxWideCards = 1000

// Database provides an interface to PostgreSQL
type Database struct {
	db *sql.DB
}

// NewDatabase returns a data Database
func NewDatabase(options ...string) (*Database, error) {
	if len(options) != 6 {
		return nil, ErrInvalidArguments
	}

	strParams := fmt.Sprintf("dbname=%s host=%s port=%s user=%s password=%s sslmode=%s",
		options[0], options[1], options[2], options[3],
		strings.Replace(options[4], "'", "\\'", -1),
		options[5])

	db, err := sql.Open("postgres", strParams)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
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

	maxChuck := maxWideCards / wcPreRecord
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
