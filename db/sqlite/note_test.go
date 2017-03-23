package sqlite

import (
	"testing"

	"github.com/anmil/quicknote/note"
	"github.com/anmil/quicknote/test"
)

func TestCreateNoteSQLiteUnit(t *testing.T) {
	db := openDatabase(t)
	defer closeDatabase(db, t)

	n := test.GetTestNotes()[0]
	saveNote(t, db, n)

	getNoteByID(t, db, n)
}

func TestGetNoteSQLiteUnit(t *testing.T) {
	db := openDatabase(t)
	defer closeDatabase(db, t)

	notes := test.GetTestNotes()
	saveNotes(t, db, notes)

	getNoteByID(t, db, notes[0])
	getNotesByID(t, db, notes)
	getNotesByBook(t, db, notes)
	getNotesAll(t, db, notes)
}

func TestEditNoteSQLiteUnit(t *testing.T) {
	db := openDatabase(t)
	defer closeDatabase(db, t)

	bk := note.NewBook()
	bk.Name = "NewBook"

	err := db.CreateBook(bk)
	if err != nil {
		t.Fatal(err)
	}

	var ids []int64
	notes := test.GetTestNotes()

	for _, n := range notes {
		saveNote(t, db, n)
		ids = append(ids, n.ID)
		n.Book = bk
	}

	if err := db.EditNoteByIDBook(ids, bk); err != nil {
		t.Fatal(err)
	}

	getNotesByBook(t, db, notes)
}

func TestEditNoteBookSQLiteUnit(t *testing.T) {
	db := openDatabase(t)
	defer closeDatabase(db, t)

	n := test.GetTestNotes()[0]
	saveNote(t, db, n)

	getNoteByID(t, db, n)

	n.Title = "New title"
	if err := db.EditNote(n); err != nil {
		t.Fatal(err)
	}

	getNoteByID(t, db, n)
}

func TestDeleteNoteSQLiteUnit(t *testing.T) {
	db := openDatabase(t)
	defer closeDatabase(db, t)

	n := test.GetTestNotes()[0]
	saveNote(t, db, n)

	getNoteByID(t, db, n)

	if err := db.DeleteNote(n); err != nil {
		t.Fatal(err)
	}

	if nn, err := db.GetNoteByID(n.ID); err != nil {
		t.Fatal(err)
	} else if nn != nil {
		t.Fatal("Expected nil, got a note")
	}
}

func saveNotes(t *testing.T, db *Database, notes []*note.Note) {
	for _, n := range notes {
		saveNote(t, db, n)
	}
}

func saveNote(t *testing.T, db *Database, n *note.Note) {
	if bk, err := db.GetBookByName(n.Book.Name); err != nil {
		t.Fatal(err)
	} else if bk == nil {
		if err := db.CreateBook(n.Book); err != nil {
			t.Fatal(err)
		}
	}

	for _, tag := range n.Tags {
		if bk, err := db.GetTagByName(tag.Name); err != nil {
			t.Fatal(err)
		} else if bk == nil {
			if err := db.CreateTag(tag); err != nil {
				t.Fatal(err)
			}
		}
	}

	if err := db.CreateNote(n); err != nil {
		t.Fatal(err)
	}
}

func getNoteByID(t *testing.T, db *Database, n *note.Note) {
	if nn, err := db.GetNoteByID(n.ID); err != nil {
		t.Fatal(err)
	} else if nn == nil {
		t.Fatal("Expected 1 note, got nil")
	} else if nn.ID != n.ID {
		t.Fatalf("Expected note with ID %d, got %d", n.ID, nn.ID)
	} else {
		test.CheckTags(t, nn.Tags, n.Tags)
	}
}

func getNotesByID(t *testing.T, db *Database, notes []*note.Note) {
	var ids []int64
	for _, n := range notes {
		ids = append(ids, n.ID)
	}

	if nn, err := db.GetAllNotesByIDs(ids); err != nil {
		t.Fatal(err)
	} else if len(nn) != len(notes) {
		t.Fatalf("Expected %d notes, got %d", len(notes), len(nn))
	} else {
		test.CheckNotes(t, nn, notes)
		for i := 0; i < len(nn); i++ {
			test.CheckTags(t, nn[i].Tags, notes[i].Tags)
		}
	}
}

func getNotesByBook(t *testing.T, db *Database, notes []*note.Note) {
	if nn, err := db.GetAllBookNotes(notes[0].Book, "modified", "asc"); err != nil {
		t.Fatal(err)
	} else if len(nn) != len(notes) {
		t.Fatalf("Expected %d notes, got %d", len(notes), len(nn))
	} else {
		test.CheckNotes(t, nn, notes)
		for i := 0; i < len(nn); i++ {
			test.CheckTags(t, nn[i].Tags, notes[i].Tags)
		}
	}
}

func getNotesAll(t *testing.T, db *Database, notes []*note.Note) {
	if nn, err := db.GetAllNotes("modified", "asc"); err != nil {
		t.Fatal(err)
	} else if len(nn) != len(notes) {
		t.Fatalf("Expected %d notes, got %d", len(notes), len(nn))
	} else {
		test.CheckNotes(t, nn, notes)
		for i := 0; i < len(nn); i++ {
			test.CheckTags(t, nn[i].Tags, notes[i].Tags)
		}
	}
}
