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

import "fmt"

var (
	Black  = 0
	Gray   = 8
	Red    = 1
	Green  = 2
	Yellow = 3
	Blue   = 4
	Pink   = 5
	Teal   = 6
	White  = 7

	HiRed    = 9
	HiGreen  = 10
	HiYellow = 11
	HiBlue   = 12
	HiPink   = 13
	HiTeal   = 14
	HiWhite  = 15
)

func colorFG(s string, fg int) string {
	if fg < 0 || fg > 255 {
		return s
	}

	return fmt.Sprintf("\x1b[38;5;%dm%s\x1b[0m", fg, s)
}

func colorBG(s string, bg int) string {
	if bg < 0 || bg > 255 {
		return s
	}

	return fmt.Sprintf("\x1b[48;5;%dm\x1b[30m%s\x1b[0m", bg, s)
}

func colorFBG(s string, fg int, bg int) string {
	return colorBG(colorFG(s, fg), bg)
}
