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
	"strconv"

	"github.com/anmil/quicknote"

	"github.com/spf13/cobra"
)

func init() {
	DeleteCmd.AddCommand(DeleteNoteCmd)
}

// DeleteNoteCmd Delete Note from Book
var DeleteNoteCmd = &cobra.Command{
	Use:   "note <note id...>",
	Short: "Delete Note from Book",
	Run:   deleteNoteCmdRun,
}

func deleteNoteCmdRun(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		for _, id := range args {
			noteID, err := strconv.ParseInt(id, 10, 64)
			exitOnError(err)

			n, err := dbConn.GetNoteByID(noteID)
			exitOnError(err)

			deleteNote(n)
		}

		fmt.Println("Note(s) deleted")
	} else {
		exitValidationError("No Note ids provided", cmd)
	}
}

func deleteNote(n *quicknote.Note) {
	if n == nil {
		return
	}

	err := dbConn.DeleteNote(n)
	exitOnError(err)

	err = idxConn.DeleteNote(n)
	exitOnError(err)
}
