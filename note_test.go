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

package quicknote

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
	n.Created = time.Unix(1490020989, 0).UTC()
	n.Modified = time.Unix(1490020989, 0).UTC()
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
	answer := `{"id":123456,"created":"2017-03-20T14:43:09Z","modified":"2017-03-20T14:43:09Z","type":"basic","title":"Json Title Test","body":"Json body test","book":"TestBook","tags":["tag1","tag2"]}`

	if results != answer {
		t.Error("JSON strings do not match")
	}
}
