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

import "github.com/jroimartin/gocui"

//================ Main events callbacks ==========================
func (c *CUI) quitCB(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

//================ Note List view events callbacks ================
func (c *CUI) moveNVSelUpCB(g *gocui.Gui, v *gocui.View) error {
	return c.NoteListView.MoveSelectionUp()
}

func (c *CUI) moveNVSelDnCB(g *gocui.Gui, v *gocui.View) error {
	return c.NoteListView.MoveSelectionDown()
}

func (c *CUI) moveNVSelPUpCB(g *gocui.Gui, v *gocui.View) error {
	return c.NoteListView.PageSelectionUp()
}

func (c *CUI) moveNVSelPDnCB(g *gocui.Gui, v *gocui.View) error {
	return c.NoteListView.PageSelectionDown()
}

func (c *CUI) changeBookCB(g *gocui.Gui, v *gocui.View) error {
	books, err := c.DBConn.GetAllBooks()
	if err != nil {
		return err
	}

	maxX, maxY := g.Size()
	bs, err := NewBookSelectorV(c, 5, 5, maxX-5, maxY-5)
	if err != nil {
		return err
	}

	bs.SetSelectedBook(c.WBook)

	if err = bs.SetBooks(books); err != nil {
		return err
	}

	c.SetCurrentView(bs)
	c.BookSelector = bs

	return nil
}

//================ Book Selector view events callbacks ============
func (c *CUI) moveBSSelUpCB(g *gocui.Gui, v *gocui.View) error {
	return c.BookSelector.MoveSelectionUp()
}

func (c *CUI) moveBSSelDnCB(g *gocui.Gui, v *gocui.View) error {
	return c.BookSelector.MoveSelectionDown()
}

func (c *CUI) moveBSSelLfCB(g *gocui.Gui, v *gocui.View) error {
	return c.BookSelector.MoveSelectionLeft()
}

func (c *CUI) moveBSSelRtCB(g *gocui.Gui, v *gocui.View) error {
	return c.BookSelector.MoveSelectionRight()
}

func (c *CUI) selBSCB(g *gocui.Gui, v *gocui.View) error {
	c.WBook = c.BookSelector.GetSelectedBook()
	if err := c.closeChangeBookCB(g, v); err != nil {
		return err
	}

	notes, err := c.DBConn.GetAllBookNotes(c.WBook, "modified", "asc")
	if err != nil {
		return err
	}

	if err = c.StatusBarView.SetWorkingBookName(c.WBook.Name); err != nil {
		return err
	}

	return c.NoteListView.SetNotes(notes)
}

func (c *CUI) closeChangeBookCB(g *gocui.Gui, v *gocui.View) error {
	if err := c.GoCUI.DeleteView(c.BookSelector.Name()); err != nil {
		return err
	}
	c.BookSelector = nil
	return c.SetCurrentView(c.NoteListView)
}
