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

func TestGetTagByNamePostgresIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestGetTagByNamePostgresIntegration in short mode")
	}

	db := openDatabase(t)
	defer closeDatabase(db, t)

	n := test.GetTestNotes()[0]
	saveNote(t, db, n)

	if tag, err := db.GetTagByName(n.Tags[0].Name); err != nil {
		t.Fatal(err)
	} else if tag == nil {
		t.Fatal("Expected 1 tag, got nil")
	} else if tag.Name != n.Tags[0].Name {
		t.Fatalf("Expected tag %s, got %s", n.Tags[0].Name, tag.Name)
	}
}

func TestGetOrCreateTagPostgresIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestGetOrCreateTagPostgresIntegration in short mode")
	}

	db := openDatabase(t)
	defer closeDatabase(db, t)

	tag1, err := db.GetOrCreateTagByName("NewTag")
	if err != nil {
		t.Fatal(err)
	} else if tag1 == nil {
		t.Fatal("Expected 1 tag, got nil")
	}

	if tag2, err := db.GetTagByName(tag1.Name); err != nil {
		t.Fatal(err)
	} else if tag2 == nil {
		t.Fatal("Expected 1 tag, got nil")
	} else if tag2.Name != tag1.Name {
		t.Fatalf("Expected tag %s, got %s", tag1.Name, tag2.Name)
	}
}

func TestGetTagsPostgresIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestGetTagsPostgresIntegration in short mode")
	}

	db := openDatabase(t)
	defer closeDatabase(db, t)

	notes := test.GetTestNotes()
	saveNotes(t, db, notes)

	tags, err := db.GetAllBookTags(notes[0].Book)
	if err != nil {
		t.Fatal(err)
	} else {
		test.CheckTags(t, test.AllTags, tags)
	}

	if tags, err = db.GetAllTags(); err != nil {
		t.Fatal(err)
	} else {
		test.CheckTags(t, test.AllTags, tags)
	}
}

func TestLoadNoteTagsPostgresIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestLoadNoteTagsPostgresIntegration in short mode")
	}

	db := openDatabase(t)
	defer closeDatabase(db, t)

	n := test.GetTestNotes()[0]
	saveNote(t, db, n)

	tags := n.Tags
	n.Tags = make(note.Tags, 0)

	if err := db.LoadNoteTags(n); err != nil {
		t.Fatal(err)
	} else {
		test.CheckTags(t, tags, n.Tags)
	}
}
