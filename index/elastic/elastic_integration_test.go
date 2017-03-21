package elastic

import (
	"fmt"
	"testing"
	"time"

	"github.com/anmil/quicknote/test"
)

var indexName = "qnote-test"
var indexHost = "http://127.0.0.1:9200"
var index *Index

func TestIndexIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestIndexNoteIntegration in short mode")
	}

	var err error
	index, err = NewIndex(indexHost, indexName)
	test.CheckErrorFatal(t, err)

	t.Run("index-note", testIndexNote)
	t.Run("index-notes", testIndexNotes)
	t.Run("search-note", testSearchNote)
	t.Run("search-phrase-note", testSearchNotePhrase)
	t.Run("delete-note", testDeleteNote)
	t.Run("delete-book", testDeleteBook)

	err = index.DeleteIndex()
	test.CheckError(t, err)
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

	// Have to wait at least 1 second for the index to complete
	// I hate this, but there is no API to block till completion
	// https://github.com/elastic/elasticsearch/issues/1063
	time.Sleep(time.Millisecond * 1000)

	query := fmt.Sprintf("id:%d", n.ID)
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

	time.Sleep(time.Millisecond * 1000)

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

	time.Sleep(time.Millisecond * 1000)

	query := fmt.Sprintf("id:%d", n.ID)
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

	time.Sleep(time.Millisecond * 1000)

	query = fmt.Sprintf("id:%d", n.ID)
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

	time.Sleep(time.Millisecond * 1000)

	query := fmt.Sprintf("book:%s", n.Book.Name)
	ids, total, err := index.SearchNote(query, 10, 0)
	test.CheckError(t, err)

	if int(total) != len(test.TestNotes) {
		t.Errorf("Expected %d results, got %d", len(test.TestNotes), total)
	} else if len(ids) != len(test.TestNotes) {
		t.Errorf("Expected %d ID, got %d", len(test.TestNotes), len(ids))
	}

	err = index.DeleteBook(n.Book)
	test.CheckError(t, err)

	time.Sleep(time.Millisecond * 1000)

	query = fmt.Sprintf("book:%s", n.Book.Name)
	_, total, err = index.SearchNote(query, 10, 0)
	test.CheckError(t, err)

	if total != 0 {
		t.Errorf("Expected 0 results, got %d", total)
	}
}
