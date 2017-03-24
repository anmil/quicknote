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
	"testing"

	"github.com/anmil/quicknote/test"
)

var tableNames = []string{
	"books",
	"note_book_tag",
	"note_tag",
	"notes",
	"sqlite_sequence",
	"tags",
}

func openDatabase(t *testing.T) *Database {
	db, err := NewDatabase("file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func closeDatabase(db *Database, t *testing.T) {
	if err := db.Close(); err != nil {
		t.Error(err)
	}
}

func TestCreateDatabaseSQLite(t *testing.T) {
	db := openDatabase(t)
	defer closeDatabase(db, t)

	if tables, err := db.GetTableNames(); err != nil {
		t.Fatal(err)
	} else if !test.StringSliceEq(tables, tableNames) {
		t.Fatal("Database either has extra or is missing tables")
	}
}
