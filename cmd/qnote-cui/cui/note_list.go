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
	"fmt"

	"github.com/anmil/quicknote/note"
	"github.com/jroimartin/gocui"
)

// NoteListVN view name
var NoteListVN = "note_list"

// NoteListV displays a list of notes allowing
// the user to scroll and page.
type NoteListV struct {
	View

	// Working Book name
	notes     []*note.Note
	highestID int64

	selIdx int

	// Generic Status Message
	msg string
}

// NewNoteListV returns a new Note List view
func NewNoteListV(c *CUI, x0, y0, x1, y1 int) (*NoteListV, error) {
	v, err := c.GoCUI.SetView(NoteListVN, x0, y0, x1, y1)
	if err != nil && err != gocui.ErrUnknownView {
		return nil, err
	}

	v.Editable = false
	v.Wrap = false
	v.Frame = false

	n := &NoteListV{
		View: View{
			c:  c,
			v:  v,
			x0: x0,
			y0: y0,
			x1: x1,
			y1: y1,
		},
	}

	return n, nil
}

// SetNotes sets the list of Notes to display and renders them
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

// MoveSelectionUp moves selection up by one
func (n *NoteListV) MoveSelectionUp() error {
	return n.SetSelection(n.selIdx - 1)
}

// MoveSelectionDown moves selection down by one
func (n *NoteListV) MoveSelectionDown() error {
	return n.SetSelection(n.selIdx + 1)
}

// HalfPageUp moves selection halfway the view's hight up
func (n *NoteListV) HalfPageUp() error {
	_, sy := n.v.Size()
	return n.SetSelection(n.selIdx - (sy / 2))
}

// HalfPageDown moves selection halfway the view's hight down
func (n *NoteListV) HalfPageDown() error {
	_, sy := n.v.Size()
	return n.SetSelection(n.selIdx + (sy / 2))
}

// PageSelectionUp moves selection the length of the view's hight up
func (n *NoteListV) PageSelectionUp() error {
	_, sy := n.v.Size()
	return n.SetSelection(n.selIdx - sy)
}

// PageSelectionDown moves selection the length of the view's hight down
func (n *NoteListV) PageSelectionDown() error {
	_, sy := n.v.Size()
	return n.SetSelection(n.selIdx + sy)
}

// SetSelection moves selection to the given index
func (n *NoteListV) SetSelection(idx int) error {
	_, sy := n.v.Size()
	cx, cy := n.v.Cursor()
	ox, oy := n.v.Origin()

	// If we get a index that is outside the bounds of the list
	// we just set it to the end or beginning of the list
	if idx >= len(n.notes) {
		idx = len(n.notes) - 1
	}
	if idx < 0 {
		idx = 0
	}

	// Calculate the new cursor and origin position using the
	// old and new index.
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

	// Once we finish the calculations
	// we can override the old one
	n.selIdx = idx

	// Set the origin first so the cursor is in the visitable view
	// otherwise we will get an error if it is not.
	if err := n.v.SetOrigin(ox, oy); err != nil {
		return err
	}
	if err := n.v.SetCursor(cx, cy); err != nil {
		return err
	}

	// I would love to not have to re-render the entire view, but
	// I have not found a way to insert text at the cursor without
	// losing color
	return n.Render()
}

// Resize sets the a new size to render the view
func (n *NoteListV) Resize(x0, y0, x1, y1 int) error {
	n.x0 = x0
	n.y0 = y0
	n.x1 = x1
	n.y1 = y1
	return n.Render()
}

// Render creates the list of Notes and renders them to the view. Highlighting
// the selected Note.
func (n *NoteListV) Render() error {
	_, err := n.c.GoCUI.SetView(n.Name(), n.x0, n.y0, n.x1, n.y1)
	if err != nil {
		return err
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
