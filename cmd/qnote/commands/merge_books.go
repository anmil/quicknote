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
	RootCmd.AddCommand(MergeBooksCmd)
}

// MergeBooksCmd Merge one book into another
var MergeBooksCmd = &cobra.Command{
	Use:   "merge [flags] <book_name 1> <book_name 2>",
	Short: "Merge all notes from <book_name 1> into <book_name 2>",
	Long: `Merge all of the notes from <book_name 1> into <book_name 2>. Than <book_name 1>
is deleted and <book_name 2> is re-indexed.`,
	Run: mergeBooksCmdRun,
}

func mergeBooksCmdRun(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		exitValidationError("two book names must be given", cmd)
	}

	book1, err := dbConn.GetBookByName(args[0])
	exitOnError(err)
	book2, err := dbConn.GetBookByName(args[1])
	exitOnError(err)

	if book1 == nil {
		exitValidationError(fmt.Sprintf("Book %s does not exists", args[0]), cmd)
	}
	if book2 == nil {
		exitValidationError(fmt.Sprintf("Book %s does not exists", args[1]), cmd)
	}

	cMsg := "This will merge all of the notes from Book %s into Book %s and than delete Book %s, are you sure?"
	if skipConfirm || utils.AskForConfirmationMust(fmt.Sprintf(cMsg, args[0], args[1], args[0])) {
		err = dbConn.MergeBooks(book1, book2)
		exitOnError(err)

		err = idxConn.DeleteBook(book1)

		notes, err := dbConn.GetAllBookNotes(book2, sortBy, displayOrder)
		exitOnError(err)

		err = idxConn.IndexNotes(notes)
		exitOnError(err)

		fmt.Println("Books merged")
	}
}
