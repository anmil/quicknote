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
	"bytes"
	"errors"
	"fmt"

	"github.com/anmil/quicknote/note"
	"github.com/jroimartin/gocui"
)

var NoteListVN = "note_list"

type NoteListV struct {
	c *CUI
	v *gocui.View

	x0 int
	y0 int
	x1 int
	y1 int

	// Working Book name
	notes     []*note.Note
	highestID int64

	selIdx int

	// Generic Status Message
	msg string
}

func NewNoteListV(c *CUI, x0, y0, x1, y1 int) (*NoteListV, error) {
	v, err := c.GoCUI.SetView(NoteListVN, x0, y0, x1, y1)
	if err != nil && err != gocui.ErrUnknownView {
		return nil, err
	}

	v.Editable = false
	v.Wrap = false
	v.Frame = false

	// v.Highlight = true
	// v.SelBgColor = gocui.ColorGreen

	n := &NoteListV{c: c, v: v}
	n.x0 = x0
	n.y0 = y0
	n.x1 = x1
	n.y1 = y1

	return n, nil
}

func (n *NoteListV) SetNotes(ns []*note.Note) error {
	n.selIdx = 0
	n.notes = ns

	for _, note := range n.notes {
		if note.ID > n.highestID {
			n.highestID = note.ID
		}
	}

	return n.Render()
}

func (n *NoteListV) MoveSelectionUp(g *gocui.Gui, v *gocui.View) error {
	return n.MoveSelection(n.selIdx - 1)
}

func (n *NoteListV) MoveSelectionDown(g *gocui.Gui, v *gocui.View) error {
	return n.MoveSelection(n.selIdx + 1)
}

func (n *NoteListV) HalfPageUp(g *gocui.Gui, v *gocui.View) error {
	_, sy := n.v.Size()
	return n.MoveSelection(n.selIdx - (sy / 2))
}

func (n *NoteListV) HalfPageDown(g *gocui.Gui, v *gocui.View) error {
	_, sy := n.v.Size()
	return n.MoveSelection(n.selIdx + (sy / 2))
}

func (n *NoteListV) PageUp(g *gocui.Gui, v *gocui.View) error {
	_, sy := n.v.Size()
	return n.MoveSelection(n.selIdx - sy)
}

func (n *NoteListV) PageDown(g *gocui.Gui, v *gocui.View) error {
	_, sy := n.v.Size()
	return n.MoveSelection(n.selIdx + sy)
}

func (n *NoteListV) MoveSelection(idx int) error {
	_, sy := n.v.Size()
	cx, cy := n.v.Cursor()
	ox, oy := n.v.Origin()

	if idx >= len(n.notes) {
		idx = len(n.notes) - 1
	}
	if idx < 0 {
		idx = 0
	}

	if idx-n.selIdx == 0 || idx < 0 {
		return nil
	} else if idx-n.selIdx >= 1 {
		d := idx - n.selIdx
		if idx < oy+sy {
			cy = cy + d
		} else {
			cy = sy - 1
			oy = idx - sy + 1
		}
	} else if idx-n.selIdx <= -1 {
		d := idx - n.selIdx
		if idx >= oy {
			cy = cy + d
		} else {
			cy = 0
			oy = idx
		}
	}

	n.c.StatusBarView.SetMessage(fmt.Sprintf("nid: %d cy: %d sy: %d oy: %d idx: %d cidx: %d d: %d",
		n.notes[idx].ID, cy, sy, oy, idx, n.selIdx, idx-n.selIdx))

	n.selIdx = idx

	if err := n.v.SetOrigin(ox, oy); err != nil {
		return err
	}
	if err := n.v.SetCursor(cx, cy); err != nil {
		return err
	}

	return n.Render()
}

func (n *NoteListV) Resize(x0, y0, x1, y1 int) error {
	n.x0 = x0
	n.y0 = y0
	n.x1 = x1
	n.y1 = y1
	return n.Render()
}

func (n *NoteListV) Render() error {
	_, err := n.c.GoCUI.SetView(NoteListVN, n.x0, n.y0, n.x1, n.y1)
	if err != nil {
		return errors.New("Failed to resize")
	}

	var buff bytes.Buffer

	idLen := len(fmt.Sprintf("%d", n.highestID))

	for idx, note := range n.notes {
		idStr := fmt.Sprintf(fmt.Sprintf("%%%dd", idLen), note.ID-1)

		var s string
		if idx == n.selIdx {
			s = colorFBG(fmt.Sprintf("%s: %s\n", idStr, note.Title), HiWhite, Teal)
		} else {
			s = fmt.Sprintf("%s: %s\n", colorFG(idStr, Green), note.Title)
		}
		buff.WriteString(s)
	}

	n.v.Clear()
	fmt.Fprint(n.v, buff.String())

	return nil
}
