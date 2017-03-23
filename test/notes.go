// +build integration
package test

import (
	"encoding/json"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/anmil/quicknote/note"
)

var notesJSON = `[
  {
    "id": 600,
    "created": "2017-03-20T10:12:42.783947469-04:00",
    "modified": "2017-03-20T10:12:42.783947542-04:00",
    "type": "basic",
    "title": "This is test 1 of the basic parser",
    "body": "#basic #test #parser\n\nLorem ipsum dolor sit amet, consectetur adipiscing elit.\nNulla tincidunt diam eu purus laoreet condimentum. Duis\ntempus, turpis vitae varius ullamcorper, sapien erat\ncursus lacus, et lacinia ligula dolor quis nibh.",
    "book": "test",
    "tags": [
      "basic",
      "test",
      "parser"
    ]
  },
  {
    "id": 601,
    "created": "2017-03-20T10:12:52.608585309-04:00",
    "modified": "2017-03-20T10:12:52.608585472-04:00",
    "type": "basic",
    "title": "This is #test 2 of the #basic #parser",
    "body": "Lorem ipsum dolor sit amet, consectetur adipiscing elit.\nNulla tincidunt diam eu purus laoreet condimentum. Duis\ntempus, turpis vitae varius ullamcorper, sapien erat\ncursus lacus, et lacinia ligula dolor quis nibh.",
    "book": "test",
    "tags": [
      "basic",
      "test",
      "parser"
    ]
  },
  {
    "id": 602,
    "created": "2017-03-20T10:25:30.570182485-04:00",
    "modified": "2017-03-20T10:25:30.570182563-04:00",
    "type": "basic",
    "title": "This is #test 2 of the #basic #parser",
    "body": "Lorem ipsum dolor sit amet, consectetur adipiscing elit.\nNulla tincidunt diam eu purus laoreet condimentum. Duis\ntempus, turpis vitae varius ullamcorper, sapien erat\ncursus lacus, et lacinia ligula dolor #quis nibh.#",
    "book": "test",
    "tags": [
      "basic",
      "test",
      "parser",
      "quis"
    ]
  }
]`

type JsonNote struct {
	ID       int64     `json:"id"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Type     string    `json:"type"`
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	Book     string    `json:"book"`
	Tags     []string  `json:"tags"`
}

func init() {
	noteBooks = make(map[string]*note.Book)
	noteTags = make(map[string]*note.Tag)
}

type Notes []*note.Note

func (n Notes) Len() int {
	return len(n)
}

func (n Notes) Less(i, j int) bool {
	return n[i].ID < n[j].ID
}

func (n Notes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func GetTestNotes() Notes {
	reader := strings.NewReader(notesJSON)
	dec := json.NewDecoder(reader)

	jsonNotes := make([]*JsonNote, 0)
	err := dec.Decode(&jsonNotes)
	if err != nil {
		panic(err)
	}

	testNotes := make([]*note.Note, len(jsonNotes))
	for idx, jn := range jsonNotes {
		tags := make([]*note.Tag, len(jn.Tags))
		for i, t := range jn.Tags {
			tags[i] = getTag(t)
		}

		n := note.NewNote()
		n.ID = jn.ID
		n.Created = jn.Created
		n.Modified = jn.Modified
		n.Type = jn.Type
		n.Title = jn.Title
		n.Body = jn.Body
		n.Book = getBook(jn.Book)
		n.Tags = tags

		testNotes[idx] = n
	}

	return testNotes
}

func CheckNotes(t *testing.T, notes1, notes2 []*note.Note) {
	nnNotes := Notes{}
	for _, t := range notes1 {
		nnNotes = append(nnNotes, t)
	}
	sort.Sort(nnNotes)

	nNotes := Notes{}
	for _, t := range notes2 {
		nNotes = append(nNotes, t)
	}
	sort.Sort(nNotes)

	if !NoteSliceEq(nnNotes, nNotes) {
		t.Fatal("Did not received the corrected Notes")
	}
}

func NoteSliceEq(a, b []*note.Note) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Title != b[i].Title {
			return false
		}
		if a[i].Body != b[i].Body {
			return false
		}
	}
	return true
}
