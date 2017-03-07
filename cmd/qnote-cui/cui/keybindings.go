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
	"github.com/jroimartin/gocui"
)

func (c *CUI) setKeybindings() {
	// Global key bindings
	must(c.GoCUI.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, c.quitCB))

	// Note list view key bindings
	must(c.GoCUI.SetKeybinding(NoteListVN, gocui.KeyEsc, gocui.ModNone, c.quitCB))
	must(c.GoCUI.SetKeybinding(NoteListVN, gocui.KeyArrowUp, gocui.ModNone, c.moveNVSelUpCB))
	must(c.GoCUI.SetKeybinding(NoteListVN, gocui.KeyArrowDown, gocui.ModNone, c.moveNVSelDnCB))
	must(c.GoCUI.SetKeybinding(NoteListVN, gocui.KeyPgup, gocui.ModNone, c.moveNVSelPUpCB))
	must(c.GoCUI.SetKeybinding(NoteListVN, gocui.KeyPgdn, gocui.ModNone, c.moveNVSelPDnCB))
}
