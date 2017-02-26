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

	"github.com/anmil/quicknote/cmd/shared/utils"
	"github.com/spf13/cobra"
)

func init() {
	DeleteCmd.AddCommand(DeleteBookCmd)
}

// DeleteBookCmd delete a book and all of it's Notes
var DeleteBookCmd = &cobra.Command{
	Use:   "book <book id>",
	Short: "Delete a Book and all of it's Notes",
	Run:   deleteNotebookCmdRun,
}

func deleteNotebookCmdRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("Please give only one Notebook name")
		return
	}

	bk, err := dbConn.GetBookByName(args[0])
	exitOnError(err)

	if bk == nil {
		fmt.Println("Notebook does not exists")
		return
	}

	cMsg := "This will delete all notes in this notebook, are you sure?"
	if skipConfirm || utils.AskForConfirmationMust(cMsg) {
		err = dbConn.DeleteBook(bk)
		exitOnError(err)

		err = idxConn.DeleteBook(bk)
		exitOnError(err)

		fmt.Println("Notebook deleted")
	}
}
