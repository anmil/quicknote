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
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/anmil/quicknote"
)

var AllTags quicknote.Tags
var noteTags map[string]*quicknote.Tag

func getTag(name string) *quicknote.Tag {
	if t, found := noteTags[name]; found {
		return t
	}

	t := quicknote.NewTag()
	t.ID = int64(len(noteTags) + 1)
	t.Created = time.Now()
	t.Modified = time.Now()
	t.Name = name
	noteTags[name] = t
	AllTags = append(AllTags, t)

	return t
}

func CheckTags(t *testing.T, tag1, tag2 quicknote.Tags) {
	nnTags := quicknote.Tags{}
	for _, t := range tag1 {
		nnTags = append(nnTags, t)
	}
	sort.Sort(nnTags)

	nTags := quicknote.Tags{}
	for _, t := range tag2 {
		nTags = append(nTags, t)
	}
	sort.Sort(nTags)

	if !TagSliceEq(nnTags, nTags) {
		t.Fatal("Did not received the corrected Tags")
	}
}

func TagSliceEq(a, b quicknote.Tags) bool {
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
			fmt.Println("Tags:", a[i].Name, b[i].Name)
			return false
		}
	}
	return true
}
