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

package cui

import (
	"fmt"

	"github.com/anmil/quicknote/cmd/shared/utils"
	"github.com/anmil/quicknote/note"
	"github.com/jroimartin/gocui"
)

// BookSelectorVN view name
var BookSelectorVN = "book_selector"

type BookSelectorV struct {
	View

	books   []*note.Book
	selBook *note.Book

	bkTable [][]string
}

func NewBookSelectorV(c *CUI, x0, y0, x1, y1 int) (*BookSelectorV, error) {
	v, err := c.GoCUI.SetView(BookSelectorVN, x0, y0, x1, y1)
	if err != nil && err != gocui.ErrUnknownView {
		return nil, err
	}

	v.Editable = false
	v.Wrap = false
	v.Frame = true
	v.Title = "Book Selector"

	b := &BookSelectorV{
		View: View{
			c:  c,
			v:  v,
			x0: x0,
			y0: y0,
			x1: x1,
			y1: y1,
		},
	}

	return b, nil
}

func (b *BookSelectorV) GetSelectedBook() *note.Book {
	return b.selBook
}

func (b *BookSelectorV) getBookByName(name string) *note.Book {
	for _, bk := range b.books {
		if bk.Name == name {
			return bk
		}
	}
	return nil
}

func (b *BookSelectorV) getCurSelIdex() (int, int) {
	for rIdx, row := range b.bkTable {
		for cIdx, cell := range row {
			if cell == b.selBook.Name {
				return rIdx, cIdx
			}
		}
	}

	return -1, -1
}

func (b *BookSelectorV) MoveSelectionUp() error {
	cx, cy := b.getCurSelIdex()
	b.c.StatusBarView.SetMessage(fmt.Sprintf("Up cx: %d cy: %d", cx, cy))
	if cy > 0 {
		b.c.StatusBarView.SetMessage(fmt.Sprintf("Up cx: %d cy: %d ncx: %d ncy: %d", cx, cy, cx, cy-1))
		b.selBook = b.getBookByName(b.bkTable[cx][cy-1])
	}
	return b.Render()
}

func (b *BookSelectorV) MoveSelectionDown() error {
	cx, cy := b.getCurSelIdex()
	b.c.StatusBarView.SetMessage(fmt.Sprintf("Dn cx: %d cy: %d", cx, cy))
	if cy < len(b.bkTable[cx])-1 {
		b.c.StatusBarView.SetMessage(fmt.Sprintf("Dn cx: %d cy: %d ncx: %d ncy: %d", cx, cy, cx, cy+1))
		b.selBook = b.getBookByName(b.bkTable[cx][cy+1])
	}
	return b.Render()
}

func (b *BookSelectorV) MoveSelectionLeft() error {
	cx, cy := b.getCurSelIdex()
	b.c.StatusBarView.SetMessage(fmt.Sprintf("Lf cx: %d cy: %d", cx, cy))
	if cx > 0 {
		b.selBook = b.getBookByName(b.bkTable[cx-1][cy])
	}
	return b.Render()
}

func (b *BookSelectorV) MoveSelectionRight() error {
	cx, cy := b.getCurSelIdex()
	b.c.StatusBarView.SetMessage(fmt.Sprintf("Lf cx: %d cy: %d", cx, cy))
	if cx < len(b.bkTable)-1 && cy < len(b.bkTable[cx+1]) {
		b.selBook = b.getBookByName(b.bkTable[cx+1][cy])
	}
	return b.Render()
}

func (b *BookSelectorV) SetSelectedBook(book *note.Book) {
	b.selBook = book
}

func (b *BookSelectorV) SetBooks(books []*note.Book) error {
	b.books = books
	return b.Render()
}

// Resize sets the a new size to render the view
func (b *BookSelectorV) Resize(x0, y0, x1, y1 int) error {
	b.x0 = x0
	b.y0 = y0
	b.x1 = x1
	b.y1 = y1
	return b.Render()
}

func (b *BookSelectorV) Render() error {
	_, err := b.c.GoCUI.SetView(b.Name(), b.x0, b.y0, b.x1, b.y1)
	if err != nil {
		return err
	}

	bkNames := make([]string, len(b.books))
	for idx, book := range b.books {
		bkNames[idx] = book.Name
	}

	sx, _ := b.v.Size()

	cb := func(w string) string {
		if w == b.selBook.Name {
			return colorFG(w, Green)
		}
		return w
	}

	var grid string
	b.bkTable, grid = utils.BuildGridStringCB(bkNames, sx, cb)

	b.v.Clear()
	fmt.Fprint(b.v, grid)
	return nil
}
