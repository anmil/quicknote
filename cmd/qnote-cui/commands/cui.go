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

package commands

import (
	"bytes"
	"fmt"

	"github.com/anmil/quicknote"
	"github.com/jroimartin/gocui"
)

var (
	curSearchResultsNotes quicknote.Notes
)

func init() {
	curSearchResultsNotes = make(quicknote.Notes, 0, 0)
}

func mainLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("search_box", 0, 0, maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = false
		v.Editor = &LiveSearchEditor{
			Gui:           g,
			InputCallback: searchBoxViewEvent,
		}

		if _, err := g.SetCurrentView("search_box"); err != nil {
			return err
		}
	}
	if v, err := g.SetView("results_list", 0, 3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.SelBgColor = gocui.ColorBlue
	}
	return nil
}

func searchBoxViewEvent(g *gocui.Gui, v *gocui.View) error {
	var query string
	var err error

	_, cy := v.Cursor()
	if query, err = v.Line(cy); err != nil {
		return err
	}

	// There is a null character at the end we must remove
	query = string(bytes.Trim([]byte(query), "\x00"))

	rV, err := g.View("results_list")
	if err != nil {
		return err
	}

	_, sy := rV.Size()
	ids, _, err := idxConn.SearchNotePhrase(query, workingNotebook, "desc", sy, 0)
	if err != nil {
		return err
	}

	var highestID int64
	for _, id := range ids {
		if id > highestID {
			highestID = id
		}
	}

	idLen := len(fmt.Sprintf("%d", highestID))

	rV.Clear()
	curSearchResultsNotes = make(quicknote.Notes, 0, len(ids))
	for _, noteID := range ids {
		n, err := dbConn.GetNoteByID(noteID)
		if err != nil {
			return err
		}

		if n != nil {
			curSearchResultsNotes = append(curSearchResultsNotes, n)
			idStr := fmt.Sprintf(fmt.Sprintf("%%%dd", idLen), n.ID)
			lines := fmt.Sprintf("\x1b[38;5;50m%s\x1b[0m: %s %s", idStr, n.Book.Name, n.Title)
			fmt.Fprintln(rV, lines)
		} else {
			fmt.Fprintln(rV, "Got Null note for", noteID)
		}
	}

	return nil
}

func searchBoxKeyUpEvent(g *gocui.Gui, v *gocui.View) error {
	rV, err := g.View("results_list")
	if err != nil {
		return err
	}

	if rV != nil {
		cx, cy := rV.Cursor()
		if cy > 0 {
			if err := rV.SetCursor(cx, cy-1); err != nil {
				ox, oy := rV.Origin()
				if err := rV.SetOrigin(ox, oy-1); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func searchBoxKeyDownEvent(g *gocui.Gui, v *gocui.View) error {
	rV, err := g.View("results_list")
	if err != nil {
		return err
	}

	if rV != nil {
		cx, cy := rV.Cursor()
		if cy < len(curSearchResultsNotes)-1 {
			if err := rV.SetCursor(cx, cy+1); err != nil {
				ox, oy := rV.Origin()
				if err := rV.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func displayNote(g *gocui.Gui, v *gocui.View) error {
	if len(curSearchResultsNotes) == 0 {
		return nil
	}

	g.Cursor = false

	rV, err := g.View("results_list")
	if err != nil {
		return err
	}

	_, cy := rV.Cursor()
	n := curSearchResultsNotes[cy]

	maxX, maxY := g.Size()
	if nv, err := g.SetView("note_display", -1, -1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		nv.Wrap = true

		fmt.Fprintf(nv, "\x1b[38;5;50mID\x1b[0m: %d \x1b[38;5;50mCreated\x1b[0m: %s \x1b[38;5;50mModified\x1b[0m: %s\n\x1b[38;5;50mTitle\x1b[0m: %s\n\n%s",
			n.ID, n.Created.Format("2006-01-02 03:04:05 PM"),
			n.Modified.Format("2006-01-02 03:04:05 PM"), n.Title, n.Body)
		if _, err = g.SetCurrentView("note_display"); err != nil {
			return err
		}
	}
	return nil
}

func delDisplayNote(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = true

	if err := g.DeleteView("note_display"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("search_box"); err != nil {
		return err
	}
	return nil
}
