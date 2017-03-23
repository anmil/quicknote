package test

import (
	"sort"
	"testing"

	"github.com/anmil/quicknote/note"
)

type Books []*note.Book

func (b Books) Len() int {
	return len(b)
}

func (b Books) Less(i, j int) bool {
	return b[i].ID < b[j].ID
}

func (b Books) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

var AllBooks Books
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

func CheckBooks(t *testing.T, bks1, bks2 Books) {
	nnBks := Books{}
	for _, t := range bks1 {
		nnBks = append(nnBks, t)
	}
	sort.Sort(nnBks)

	nBks := Books{}
	for _, t := range bks2 {
		nBks = append(nBks, t)
	}
	sort.Sort(nBks)

	if !BookSliceEq(nnBks, nBks) {
		t.Fatal("Did not received the corrected Books")
	}
}

func BookSliceEq(a, b Books) bool {
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
