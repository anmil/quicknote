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
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	EditCmd.AddCommand(EditBookCmd)
}

// EditBookCmd edit Book's name
var EditBookCmd = &cobra.Command{
	Use:   "book <new book_name>",
	Short: "Edit working Book's name",
	Long:  `Edit the working Book's name. This requires re-index the Book`,
	Run:   editBookCmdRun,
}

func editBookCmdRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		exitValidationError("No name given", cmd)
	}

	workingNotebook.Name = args[0]
	err := dbConn.EditBook(workingNotebook)
	exitOnError(err)

	notes, err := dbConn.GetAllBookNotes(workingNotebook, sortBy, displayOrder)
	exitOnError(err)

	err = idxConn.IndexNotes(notes)
	exitOnError(err)

	fmt.Println("Book's name changed")
}
