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
	"github.com/anmil/quicknote/db"
	"github.com/anmil/quicknote/index"
	"github.com/anmil/quicknote/note"
	"github.com/jroimartin/gocui"
)

type CUI struct {
	GoCUI *gocui.Gui

	WBook   *note.Book
	DBCoon  db.DB
	IdxConn index.Index

	StatusBarView *StatusBarV
	NoteListView  *NoteListV
}

func NewCUI(wBook *note.Book, dbCoon db.DB, idxConn index.Index) (*CUI, error) {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		return nil, err
	}
	g.InputEsc = true

	c := &CUI{GoCUI: g, WBook: wBook, DBCoon: dbConn, IdxConn: idxConn}
	g.SetManager(c)

	maxX, maxY := g.Size()
	sb, err := NewStatusBarV(c, -1, maxY-2, maxX, maxY)
	if err != nil {
		return nil, err
	}
	if err = sb.SetWorkingBookName(c.WBook.Name); err != nil {
		return nil, err
	}
	if err = sb.SetMessage("This is a test This is a test"); err != nil {
		return nil, err
	}
	c.StatusBarView = sb

	nl, err := NewNoteListV(c, -1, -1, maxX, maxY-1)
	if err != nil {
		return nil, err
	}
	notes, err := dbCoon.GetAllBookNotes(wBook, "modified", "asc")
	if err != nil {
		return nil, err
	}
	if err = nl.SetNotes(notes); err != nil {
		return nil, err
	}
	c.NoteListView = nl

	c.setKeybindings()

	if _, err = g.SetCurrentView(NoteListVN); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *CUI) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	c.StatusBarView.Resize(-1, maxY-2, maxX, maxY)
	c.NoteListView.Resize(-1, -1, maxX, maxY-1)
	return nil
}

func (c *CUI) Cursor(b bool) {
	c.GoCUI.Cursor = b
}

func (c *CUI) InputEsc(b bool) {
	c.GoCUI.InputEsc = b
}

func (c *CUI) Run() error {
	if err := c.GoCUI.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

func (c *CUI) Close() {
	c.GoCUI.Close()
}

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

func (c *CUI) quitCB(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
