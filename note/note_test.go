package note

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNoteJSONUnit(t *testing.T) {
	bk := NewBook()
	bk.Name = "TestBook"

	tag1 := NewTag()
	tag1.Name = "tag1"

	tag2 := NewTag()
	tag2.Name = "tag2"

	n := NewNote()
	n.ID = 123456
	n.Created = time.Unix(1490020989, 0)
	n.Modified = time.Unix(1490020989, 0)
	n.Type = "basic"
	n.Title = "Json Title Test"
	n.Body = "Json body test"
	n.Book = bk
	n.Tags = append(n.Tags, tag1)
	n.Tags = append(n.Tags, tag2)

	b, err := json.Marshal(n)
	if err != nil {
		t.Error("JSON Marshal failed")
	}

	results := string(b)
	answer := `{"id":123456,"created":"2017-03-20T10:43:09-04:00","modified":"2017-03-20T10:43:09-04:00","type":"basic","title":"Json Title Test","body":"Json body test","book":"TestBook","tags":["tag1","tag2"]}`

	if results != answer {
		t.Error("JSON strings do not match")
	}
}
