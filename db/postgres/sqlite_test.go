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
	"os"
	"sort"
	"testing"

	"github.com/anmil/quicknote/test"
)

var tableNames = []string{
	"books",
	"note_book_tag",
	"note_tag",
	"notes",
	"tags",
}

var dBName string
var dBHost string
var dBPort string
var dBUser string
var dBPass string
var dBSSL string

func init() {
	dBName = os.Getenv("QN_TEST_PG_NAME")
	dBHost = os.Getenv("QN_TEST_PG_HOST")
	dBPort = os.Getenv("QN_TEST_PG_PORT")
	dBUser = os.Getenv("QN_TEST_PG_USER")
	dBPass = os.Getenv("QN_TEST_PG_PASS")
	dBSSL = os.Getenv("QN_TEST_PG_SSL")
}

func openDatabase(t *testing.T) *Database {
	options := []string{
		dBName,
		dBHost,
		dBPort,
		dBUser,
		dBPass,
		dBSSL,
	}

	db, err := NewDatabase(options...)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.ResetTables(); err != nil {
		t.Fatal(err)
	}

	return db
}

func closeDatabase(db *Database, t *testing.T) {
	if err := db.Close(); err != nil {
		t.Error(err)
	}
}

func TestCreateDatabasePostgresIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestCreateDatabasePostgresIntegration in short mode")
	}

	db := openDatabase(t)
	defer closeDatabase(db, t)

	tables, err := db.GetTableNames()

	sort.Strings(tables)
	sort.Strings(tableNames)

	if err != nil {
		t.Fatal(err)
	} else if !test.StringSliceEq(tables, tableNames) {
		t.Fatal("Database either has extra or is missing tables")
	}
}
