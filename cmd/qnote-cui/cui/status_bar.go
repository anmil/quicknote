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
	"strings"

	"github.com/jroimartin/gocui"
)

var StatusBarVN = "statusbar"

type StatusBarV struct {
	c *CUI
	v *gocui.View

	x0 int
	y0 int
	x1 int
	y1 int

	// Working Book name
	bkName string

	// Generic Status Message
	msg string
}

func NewStatusBarV(c *CUI, x0, y0, x1, y1 int) (*StatusBarV, error) {
	v, err := c.GoCUI.SetView(StatusBarVN, x0, y0, x1, y1)
	if err != nil && err != gocui.ErrUnknownView {
		return nil, err
	}

	v.Editable = false
	v.Wrap = false
	v.Frame = false

	s := &StatusBarV{c: c, v: v}
	s.x0 = x0
	s.y0 = y0
	s.x1 = x1
	s.y1 = y1

	return s, nil
}

func (s *StatusBarV) SetMessage(msg string) error {
	s.msg = msg
	return s.Render()
}

func (s *StatusBarV) SetWorkingBookName(bkName string) error {
	s.bkName = bkName
	return s.Render()
}

func (s *StatusBarV) Resize(x0, y0, x1, y1 int) error {
	s.x0 = x0
	s.y0 = y0
	s.x1 = x1
	s.y1 = y1
	return s.Render()
}

func (s *StatusBarV) Render() error {
	_, err := s.c.GoCUI.SetView(StatusBarVN, s.x0, s.y0, s.x1, s.y1)
	if err != nil {
		return err
	}

	bkStr := fmt.Sprintf("BK: %s ", s.bkName)
	rightStatus := fmt.Sprintf("  %s", bkStr)

	x, _ := s.v.Size()
	rlen := len(rightStatus)
	lMaxLen := x - rlen

	msg := s.msg
	if len(msg) > lMaxLen && len(msg) > 3 {
		msg = fmt.Sprintf("%s...", string(msg[:lMaxLen-3]))
	}

	leftStatus := fmt.Sprintf(" %s", msg)

	plen := x - rlen - len(leftStatus)
	var padding string
	if plen > 0 {
		padding = strings.Repeat(" ", plen)
	}

	statusMsg := fmt.Sprintf("%s%s%s", leftStatus, padding, rightStatus)
	statusMsg = colorFBG(statusMsg, HiWhite, Blue)

	s.v.Clear()
	fmt.Fprint(s.v, statusMsg)

	return nil
}
