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

package note

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Note types
var (
	Basic = "basic"
	URL   = "url"

	NoteTypes = []string{
		Basic,
		URL,
	}
)

// Note is our main struct for storing
// notes and their meta data.
type Note struct {
	ID       int64
	Created  time.Time
	Modified time.Time

	Type  string
	Title string
	Body  string

	Book *Book
	Tags []*Tag
}

// NewNote returns a new Note
func NewNote() *Note {
	return &Note{}
}

func (n *Note) String() string {
	bk := ""
	if n.Book != nil {
		bk = n.Book.Name
	}
	tags := ""
	if len(n.Tags) > 0 {
		tags = strings.Join(n.GetTagStringArray(), ", ")
	}
	return fmt.Sprintf("<Note ID: %d Title: %s Book: %s Tags: %s>", n.ID, n.Title, bk, tags)
}

// GetTagStringArray returns a list of the note's tag names
func (n *Note) GetTagStringArray() []string {
	tags := make([]string, 0, len(n.Tags))
	for _, tag := range n.Tags {
		tags = append(tags, tag.Name)
	}
	return tags
}

// MarshalJSON customer json Marshaler
func (n *Note) MarshalJSON() ([]byte, error) {
	tags := make([]string, len(n.Tags))
	for idx, tag := range n.Tags {
		tags[idx] = tag.Name
	}

	return json.Marshal(&struct {
		ID       int64     `json:"id"`
		Created  time.Time `json:"created"`
		Modified time.Time `json:"modified"`
		Type     string    `json:"type"`
		Title    string    `json:"title"`
		Body     string    `json:"body"`
		Book     string    `json:"book"`
		Tags     []string  `json:"tags"`
	}{
		ID:       n.ID,
		Created:  n.Created,
		Modified: n.Modified,
		Type:     n.Type,
		Title:    n.Title,
		Body:     n.Body,
		Book:     n.Book.Name,
		Tags:     tags,
	})
}
