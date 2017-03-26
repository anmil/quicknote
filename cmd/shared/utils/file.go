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

package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// GetFileExt returns the file's extension in lower case
func GetFileExt(f string) string {
	ext := filepath.Ext(f)
	return strings.ToLower(ext)
}

// ExpandFilePath takes a file path either containing environment
// variables and/or relative paths and expands it to the full path.
func ExpandFilePath(p string) (xp string, err error) {
	xp = os.ExpandEnv(p)
	if !filepath.IsAbs(xp) {
		xp, err = filepath.Abs(xp)
		if err != nil {
			return "", err
		}
	}
	return filepath.Clean(xp), err
}
