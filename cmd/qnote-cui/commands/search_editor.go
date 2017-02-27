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

import "github.com/jroimartin/gocui"

// LiveSearchEditor searches for notes as the user types
type LiveSearchEditor struct {
	Gui           *gocui.Gui
	InputCallback func(g *gocui.Gui, v *gocui.View) error
}

// Edit callback for edit events
func (ve *LiveSearchEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
		ve.InputCallback(ve.Gui, v)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
		ve.InputCallback(ve.Gui, v)
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
		ve.InputCallback(ve.Gui, v)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
		ve.InputCallback(ve.Gui, v)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	}
}
