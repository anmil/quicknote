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

func SetKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quitCB); err != nil {
		return err
	}
	if err := g.SetKeybinding("search_box", gocui.KeyEsc, gocui.ModNone, quitCB); err != nil {
		return err
	}
	if err := g.SetKeybinding("search_box", gocui.KeyEnter, gocui.ModNone, displayNote); err != nil {
		return err
	}
	if err := g.SetKeybinding("search_box", gocui.KeyArrowUp, gocui.ModNone, searchBoxKeyUpEvent); err != nil {
		return err
	}
	if err := g.SetKeybinding("search_box", gocui.KeyArrowDown, gocui.ModNone, searchBoxKeyDownEvent); err != nil {
		return err
	}
	if err := g.SetKeybinding("note_display", gocui.KeyEsc, gocui.ModNone, delDisplayNote); err != nil {
		return err
	}
	return nil
}
