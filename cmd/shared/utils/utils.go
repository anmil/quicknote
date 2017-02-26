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
	"path"
	"runtime"
)

// GetDataDirectory returns the default directory for storing user data
func GetDataDirectory() string {
	switch runtime.GOOS {
	case "linux":
		return path.Join(os.Getenv("HOME"), ".config", "quicknote")
	case "darwin":
		return path.Join(os.Getenv("HOME"), "Library", "Application Support", "quicknote")
	default:
		return path.Join(os.Getenv("HOME"), ".quicknote")
	}
}

// EnsureDirectoryExists ensures that path exists
func EnsureDirectoryExists(path string) error {
	return os.MkdirAll(path, 0700)
}

// InSliceString returns true if string i is in slice s
func InSliceString(i string, s []string) bool {
	for _, si := range s {
		if i == si {
			return true
		}
	}

	return false
}
