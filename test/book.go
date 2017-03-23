package test

import (
	"sort"
	"testing"

	"github.com/anmil/quicknote/note"
)

var AllBooks note.Books
var noteBooks map[string]*note.Book

func getBook(name string) *note.Book {
	if bk, found := noteBooks[name]; found {
		return bk
	}

	bk := note.NewBook()
	bk.Name = name
	noteBooks[name] = bk
	AllBooks = append(AllBooks, bk)

	return bk
}

func CheckBooks(t *testing.T, bks1, bks2 note.Books) {
	nnBks := note.Books{}
	for _, t := range bks1 {
		nnBks = append(nnBks, t)
	}
	sort.Sort(nnBks)

	nBks := note.Books{}
	for _, t := range bks2 {
		nBks = append(nBks, t)
	}
	sort.Sort(nBks)

	if !BookSliceEq(nnBks, nBks) {
		t.Fatal("Did not received the corrected Books")
	}
}

func BookSliceEq(a, b note.Books) bool {
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
