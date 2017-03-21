// +build integration
package test

import (
	"encoding/json"
	"strings"
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

var TestNotes []*note.Note

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

var noteBooks map[string]*note.Book

func getBook(name string) *note.Book {
	if bk, found := noteBooks[name]; found {
		return bk
	}
	bk := note.NewBook()
	bk.Name = name
	noteBooks[name] = bk
	return bk
}

var noteTags map[string]*note.Tag

func getTag(name string) *note.Tag {
	if t, found := noteTags[name]; found {
		return t
	}
	t := note.NewTag()
	t.Name = name
	noteTags[name] = t
	return t
}

func init() {
	noteBooks = make(map[string]*note.Book)
	noteTags = make(map[string]*note.Tag)

	reader := strings.NewReader(notesJSON)
	dec := json.NewDecoder(reader)

	jsonNotes := make([]*JsonNote, 0)
	err := dec.Decode(&jsonNotes)
	if err != nil {
		panic(err)
	}

	TestNotes = make([]*note.Note, len(jsonNotes))
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

		TestNotes[idx] = n
	}
}
