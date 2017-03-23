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
	"testing"

	"github.com/anmil/quicknote/note"
	"github.com/anmil/quicknote/test"
)

func TestCreateBookPostgresIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestCreateBookPostgresIntegration in short mode")
	}

	db := openDatabase(t)
	defer closeDatabase(db, t)

	bk1 := note.NewBook()
	bk1.Name = "NewBook"

	if err := db.CreateBook(bk1); err != nil {
		t.Fatal(err)
	}
}

func TestLoadBookPostgresIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestLoadBookPostgresIntegration in short mode")
	}

	db := openDatabase(t)
	defer closeDatabase(db, t)

	bk1 := note.NewBook()
	bk1.Name = "NewBook"

	if err := db.CreateBook(bk1); err != nil {
		t.Fatal(err)
	}

	bk2 := note.NewBook()
	bk2.ID = bk1.ID

	if err := db.LoadBook(bk2); err != nil {
		t.Fatal(err)
	} else if bk2.Name != bk1.Name {
		t.Fatalf("Expected book %s, got %s", bk1.Name, bk2.Name)
	}
}

func TestEditBookPostgresIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestEditBookPostgresIntegration in short mode")
	}

	db := openDatabase(t)
	defer closeDatabase(db, t)

	bk1 := note.NewBook()
	bk1.Name = "NewBook"

	if err := db.CreateBook(bk1); err != nil {
		t.Fatal(err)
	}

	bk1.Name = "EditBook"
	if err := db.EditBook(bk1); err != nil {
		t.Fatal(err)
	}

	if bk2, err := db.GetBookByName(bk1.Name); err != nil {
		t.Fatal(err)
	} else if bk2.ID != bk1.ID {
		t.Fatalf("Expected book %d, got %d", bk1.ID, bk2.ID)
	}
}

func TestGetBookByNamePostgresIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestGetBookByNamePostgresIntegration in short mode")
	}

	db := openDatabase(t)
	defer closeDatabase(db, t)

	n := test.GetTestNotes()[0]
	saveNote(t, db, n)

	if bk, err := db.GetBookByName(n.Book.Name); err != nil {
		t.Fatal(err)
	} else if bk.Name != n.Book.Name {
		t.Fatalf("Expected book %s, got %s", bk.Name, n.Book.Name)
	}
}

func TestGetOrCreateBookPostgresIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestGetOrCreateBookPostgresIntegration in short mode")
	}

	db := openDatabase(t)
	defer closeDatabase(db, t)

	bk1, err := db.GetOrCreateBookByName("NewBook")
	if err != nil {
		t.Fatal(err)
	} else if bk1 == nil {
		t.Fatal("Expected 1 book, got nil")
	}

	bk2, err := db.GetBookByName(bk1.Name)
	if err != nil {
		t.Fatal(err)
	} else if bk2 == nil {
		t.Fatal("Expected 1 book, got nil")
	} else if bk2.Name != bk1.Name {
		t.Fatalf("Expected book %s, got %s", bk1.Name, bk2.Name)
	}
}

func TestGetAllBooksPostgresIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestGetAllBooksPostgresIntegration in short mode")
	}

	db := openDatabase(t)
	defer closeDatabase(db, t)

	saveNotes(t, db, test.GetTestNotes())

	if books, err := db.GetAllBooks(); err != nil {
		t.Fatal(err)
	} else if len(books) != len(test.AllBooks) {
		t.Fatalf("Expected %d books, got %d", len(test.AllBooks), len(books))
	} else {
		test.CheckBooks(t, books, test.AllBooks)
	}
}

func TestMergeBookPostgresIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestMergeBookPostgresIntegration in short mode")
	}

	db := openDatabase(t)
	defer closeDatabase(db, t)

	notes := test.GetTestNotes()
	saveNotes(t, db, notes)

	bk1 := note.NewBook()
	bk1.Name = "NewBook"

	if err := db.CreateBook(bk1); err != nil {
		t.Fatal(err)
	}

	if err := db.MergeBooks(notes[0].Book, bk1); err != nil {
		t.Fatal(err)
	}

	for _, n := range notes {
		n.Book = bk1
	}

	getNotesByBook(t, db, notes)
}

func TestDeleteBookPostgresIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestDeleteBookPostgresIntegration in short mode")
	}

	db := openDatabase(t)
	defer closeDatabase(db, t)

	bk1 := note.NewBook()
	bk1.Name = "NewBook"

	if err := db.CreateBook(bk1); err != nil {
		t.Fatal(err)
	}

	if bk2, err := db.GetBookByName(bk1.Name); err != nil {
		t.Fatal(err)
	} else if bk2 == nil {
		t.Fatal("Expected book, got nil")
	}

	if err := db.DeleteBook(bk1); err != nil {
		t.Fatal(err)
	}

	if bk2, err := db.GetBookByName(bk1.Name); err != nil {
		t.Fatal(err)
	} else if bk2 != nil {
		t.Fatal("Expected nil, got book")
	}
}
