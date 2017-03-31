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
	"time"

	"github.com/anmil/quicknote"
)

var AllBooks quicknote.Books
var noteBooks map[string]*quicknote.Book

func getBook(name string) *quicknote.Book {
	if bk, found := noteBooks[name]; found {
		return bk
	}

	bk := quicknote.NewBook()
	bk.ID = int64(len(noteBooks) + 1)
	bk.Created = time.Now()
	bk.Modified = time.Now()
	bk.Name = name
	noteBooks[name] = bk
	AllBooks = append(AllBooks, bk)

	return bk
}

func CheckBooks(t *testing.T, bks1, bks2 quicknote.Books) {
	nnBks := quicknote.Books{}
	for _, t := range bks1 {
		nnBks = append(nnBks, t)
	}
	sort.Sort(nnBks)

	nBks := quicknote.Books{}
	for _, t := range bks2 {
		nBks = append(nBks, t)
	}
	sort.Sort(nBks)

	if !BookSliceEq(nnBks, nBks) {
		t.Fatal("Did not received the corrected Books")
	}
}

func BookSliceEq(a, b quicknote.Books) bool {
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
