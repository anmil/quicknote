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
	"os"
	"path"
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

var tempDir = path.Join(os.TempDir(), "qnote-db")

func openDatabase(t *testing.T) *Database {
	// Ensure there is no left overs
	if err := os.MkdirAll(tempDir, 0700); err != nil {
		t.Fatal(err)
	}

	file := path.Join(tempDir, "qnote.db")
	if _, err := os.Stat(file); err == nil {
		if err := os.Remove(file); err != nil {
			t.Fatal(err)
		}
	}

	// Would be nice to be able to use a memory only db. Due, to
	// the way Go sql.DB does its connection pool we can not.
	// https://groups.google.com/forum/#!msg/golang-nuts/AYZl1lNxCfA/LOr30uKy7-oJ
	db, err := NewDatabase(file)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func closeDatabase(db *Database, t *testing.T) {
	if err := db.Close(); err != nil {
		t.Error(err)
	}
	if err := os.Remove(db.DBPath); err != nil {
		t.Fatal(err)
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
