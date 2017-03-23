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

package test

import (
	"sort"
	"testing"

	"github.com/anmil/quicknote/note"
)

var AllTags note.Tags
var noteTags map[string]*note.Tag

func getTag(name string) *note.Tag {
	if t, found := noteTags[name]; found {
		return t
	}

	t := note.NewTag()
	t.Name = name
	noteTags[name] = t
	AllTags = append(AllTags, t)

	return t
}

func CheckTags(t *testing.T, tag1, tag2 note.Tags) {
	nnTags := note.Tags{}
	for _, t := range tag1 {
		nnTags = append(nnTags, t)
	}
	sort.Sort(nnTags)

	nTags := note.Tags{}
	for _, t := range tag2 {
		nTags = append(nTags, t)
	}
	sort.Sort(nTags)

	if !TagSliceEq(nnTags, nTags) {
		t.Fatal("Did not received the corrected Tags")
	}
}

func TagSliceEq(a, b note.Tags) bool {
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
		if a[i].Name != b[i].Name {
			return false
		}
	}
	return true
}
