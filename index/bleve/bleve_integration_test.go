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

func TestIndexNoteBleveIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestIndexNoteBleveIntegration in short mode")
	}

	// Ensure there is no left overs
	os.RemoveAll(tempDir)

	var err error
	index, err = NewIndex(tempDir, shardCnt)
	test.CheckErrorFatal(t, err)

	t.Run("bleve-index-note", testIndexNote)
	t.Run("bleve-index-notes", testIndexNotes)
	t.Run("bleve-search-note", testSearchNote)
	t.Run("bleve-search-phrase-note", testSearchNotePhrase)
	t.Run("bleve-delete-note", testDeleteNote)
	t.Run("bleve-delete-book", testDeleteBook)

	os.RemoveAll(tempDir)
}

func testIndexNote(t *testing.T) {
	n := test.TestNotes[0]
	err := index.IndexNote(n)
	test.CheckError(t, err)
}

func testIndexNotes(t *testing.T) {
	err := index.IndexNotes(test.TestNotes)
	test.CheckError(t, err)
}

func testSearchNote(t *testing.T) {
	n := test.TestNotes[0]
	err := index.IndexNote(n)
	test.CheckError(t, err)

	query := fmt.Sprintf("+id:%d", n.ID)
	ids, total, err := index.SearchNote(query, 10, 0)
	test.CheckError(t, err)

	if total != 1 {
		t.Errorf("Expected 1 results, got %d", total)
	} else if len(ids) != 1 {
		t.Errorf("Expected 1 ID, got %d", len(ids))
	} else if ids[0] != n.ID {
		t.Errorf("Expected ID %d, got %d", n.ID, ids[0])
	}
}

func testSearchNotePhrase(t *testing.T) {
	n := test.TestNotes[0]
	err := index.IndexNote(n)
	test.CheckError(t, err)

	query := "This is test 1 of the basic par"
	ids, total, err := index.SearchNotePhrase(query, nil, "asc", 10, 0)
	test.CheckError(t, err)

	if total != 1 {
		t.Errorf("Expected 1 results, got %d", total)
	} else if len(ids) != 1 {
		t.Errorf("Expected 1 ID, got %d", len(ids))
	} else if ids[0] != n.ID {
		t.Errorf("Expected ID %d, got %d", n.ID, ids[0])
	}
}

func testDeleteNote(t *testing.T) {
	n := test.TestNotes[0]
	err := index.IndexNote(n)
	test.CheckError(t, err)

	query := fmt.Sprintf("+id:%d", n.ID)
	ids, total, err := index.SearchNote(query, 10, 0)
	test.CheckError(t, err)

	if total != 1 {
		t.Errorf("Expected 1 results, got %d", total)
	} else if len(ids) != 1 {
		t.Errorf("Expected 1 ID, got %d", len(ids))
	} else if ids[0] != n.ID {
		t.Errorf("Expected ID %d, got %d", n.ID, ids[0])
	}

	err = index.DeleteNote(n)
	test.CheckError(t, err)

	query = fmt.Sprintf("+id:%d", n.ID)
	_, total, err = index.SearchNote(query, 10, 0)
	test.CheckError(t, err)

	if total != 0 {
		t.Errorf("Expected 0 results, got %d", total)
	}
}

func testDeleteBook(t *testing.T) {
	n := test.TestNotes[0]

	err := index.IndexNotes(test.TestNotes)
	test.CheckError(t, err)

	query := fmt.Sprintf("+book:%s", n.Book.Name)
	ids, total, err := index.SearchNote(query, 10, 0)
	test.CheckError(t, err)

	if int(total) != len(test.TestNotes) {
		t.Errorf("Expected %d results, got %d", len(test.TestNotes), total)
	} else if len(ids) != len(test.TestNotes) {
		t.Errorf("Expected %d ID, got %d", len(test.TestNotes), len(ids))
	}

	err = index.DeleteBook(n.Book)
	test.CheckError(t, err)

	query = fmt.Sprintf("book:%s", n.Book.Name)
	_, total, err = index.SearchNote(query, 10, 0)
	test.CheckError(t, err)

	if total != 0 {
		t.Errorf("Expected 0 results, got %d", total)
	}
}
