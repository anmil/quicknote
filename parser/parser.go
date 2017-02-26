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

package parser

import "errors"

// ErrParserNotSupported an unknown parser type was given
var ErrParserNotSupported = errors.New("Unsupported parser")

// Parser interface for a note parser
type Parser interface {
	Parse(text string)
	Title() string
	Tags() []string
	Body() string
}

// NewParser returns a new parser for the type given
func NewParser(ptype string) (Parser, error) {
	switch ptype {
	case "basic":
		return &BasicParser{}, nil
	default:
		return nil, ErrParserNotSupported
	}
}
