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

package bleve

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/anmil/quicknote/test"
)

var index *Index
var tempDir = path.Join(os.TempDir(), "qnote-test")
var shardCnt = 3

func TestIndexNoteBleveUnit(t *testing.T) {
	// Ensure there is no left overs
	os.RemoveAll(tempDir)
	defer os.RemoveAll(tempDir)

	var err error
	index, err = NewIndex(tempDir, shardCnt)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("bleve-index-note", testIndexNote)
	t.Run("bleve-index-notes", testIndexNotes)
	t.Run("bleve-search-note", testSearchNote)
	t.Run("bleve-search-phrase-note", testSearchNotePhrase)
	t.Run("bleve-delete-note", testDeleteNote)
	t.Run("bleve-delete-book", testDeleteBook)
}

func testIndexNote(t *testing.T) {
	n := test.GetTestNotes()[0]
	if err := index.IndexNote(n); err != nil {
		t.Fatal(err)
	}
}

func testIndexNotes(t *testing.T) {
	notes := test.GetTestNotes()
	if err := index.IndexNotes(notes); err != nil {
		t.Fatal(err)
	}
}

func testSearchNote(t *testing.T) {
	n := test.GetTestNotes()[0]
	if err := index.IndexNote(n); err != nil {
		t.Fatal(err)
	}

	query := fmt.Sprintf("+id:%d", n.ID)
	if ids, total, err := index.SearchNote(query, 10, 0); err != nil {
		t.Fatal(err)
	} else if total != 1 {
		t.Fatalf("Expected 1 results, got %d", total)
	} else if len(ids) != 1 {
		t.Fatalf("Expected 1 ID, got %d", len(ids))
	} else if ids[0] != n.ID {
		t.Fatalf("Expected ID %d, got %d", n.ID, ids[0])
	}
}

func testSearchNotePhrase(t *testing.T) {
	n := test.GetTestNotes()[0]
	if err := index.IndexNote(n); err != nil {
		t.Fatal(err)
	}

	query := "This is test 1 of the basic par"
	if ids, total, err := index.SearchNotePhrase(query, nil, "asc", 10, 0); err != nil {
		t.Fatal(err)
	} else if total != 1 {
		t.Fatalf("Expected 1 results, got %d", total)
	} else if len(ids) != 1 {
		t.Fatalf("Expected 1 ID, got %d", len(ids))
	} else if ids[0] != n.ID {
		t.Fatalf("Expected ID %d, got %d", n.ID, ids[0])
	}
}

func testDeleteNote(t *testing.T) {
	n := test.GetTestNotes()[0]
	if err := index.IndexNote(n); err != nil {
		t.Fatal(err)
	}

	query := fmt.Sprintf("+id:%d", n.ID)
	if ids, total, err := index.SearchNote(query, 10, 0); err != nil {
		t.Fatal(err)
	} else if total != 1 {
		t.Fatalf("Expected 1 results, got %d", total)
	} else if len(ids) != 1 {
		t.Fatalf("Expected 1 ID, got %d", len(ids))
	} else if ids[0] != n.ID {
		t.Fatalf("Expected ID %d, got %d", n.ID, ids[0])
	}

	if err := index.DeleteNote(n); err != nil {
		t.Fatal(err)
	}

	query = fmt.Sprintf("+id:%d", n.ID)
	if _, total, err := index.SearchNote(query, 10, 0); err != nil {
		t.Fatal(err)
	} else if total != 0 {
		t.Fatalf("Expected 0 results, got %d", total)
	}
}

func testDeleteBook(t *testing.T) {
	notes := test.GetTestNotes()
	n := notes[0]

	if err := index.IndexNotes(notes); err != nil {
		t.Fatal(err)
	}

	query := fmt.Sprintf("+book:%s", n.Book.Name)
	if ids, total, err := index.SearchNote(query, 10, 0); err != nil {
		t.Fatal(err)
	} else if int(total) != len(notes) {
		t.Fatalf("Expected %d results, got %d", len(notes), total)
	} else if len(ids) != len(notes) {
		t.Fatalf("Expected %d ID, got %d", len(notes), len(ids))
	}

	if err := index.DeleteBook(n.Book); err != nil {
		t.Fatal(err)
	}

	query = fmt.Sprintf("book:%s", n.Book.Name)
	if _, total, err := index.SearchNote(query, 10, 0); err != nil {
		t.Fatal(err)
	} else if total != 0 {
		t.Fatalf("Expected 0 results, got %d", total)
	}
}
