package sqlite

import (
	"testing"

	"github.com/anmil/quicknote/note"
	"github.com/anmil/quicknote/test"
)

func TestCreateBookSQLiteUnit(t *testing.T) {
	db := openDatabase(t)
	defer closeDatabase(db, t)

	bk1 := note.NewBook()
	bk1.Name = "NewBook"

	if err := db.CreateBook(bk1); err != nil {
		t.Fatal(err)
	}
}

func TestLoadBookSQLiteUnit(t *testing.T) {
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

func TestEditBookSQLiteUnit(t *testing.T) {
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

func TestGetBookByNameSQLiteUnit(t *testing.T) {
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

func TestGetOrCreateBookSQLiteUnit(t *testing.T) {
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

func TestGetAllBooksSQLiteUnit(t *testing.T) {
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

func TestMergeBookSQLiteUnit(t *testing.T) {
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

func TestDeleteBookSQLiteUnit(t *testing.T) {
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
