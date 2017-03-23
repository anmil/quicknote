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

package note

import (
	"fmt"
	"time"
)

// Tag is a term used as meta data for more
// accurate searching and labeling.
type Tag struct {
	ID       int64
	Created  time.Time
	Modified time.Time

	Name string
}

// NewTag returns a new Tag
func NewTag() *Tag {
	return &Tag{}
}

func (t *Tag) String() string {
	return fmt.Sprintf("<Tag ID: %d Name: %s>", t.ID, t.Name)
}
