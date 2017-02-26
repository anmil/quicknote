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

	"github.com/anmil/quicknote/cmd/shared/config"
	"github.com/anmil/quicknote/cmd/shared/utils"
	"github.com/spf13/cobra"
)

func init() {
	SplitCmd.AddCommand(SplitBookQueryCmd)
	SplitCmd.AddCommand(SplitBookIDsCmd)
}

// SplitBookQueryCmd splits one book into two book using QueryStringQuery
var SplitBookQueryCmd = &cobra.Command{
	Use:   "query [flags] <book_name> <query_string_query>",
	Short: "splits the working Book into two Books using QueryStringQuery",
	Long: `Splits the working Book into two Books. All notes matching the query will be
moved into the Book <book_name>. If <book_name> already exists, the Notes
matching the query are merged into the exciting Book. For docs on the syntax for
QueryStringQuery see the docs for 'qnote search'.`,
	Run: splitBooksQueryCmdRun,
}

func splitBooksQueryCmdRun(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		exitValidationError("invalid arguments given", cmd)
	}

	bk1 := workingNotebook

	var query string
	switch config.IndexProvider {
	case "bleve":
		query = fmt.Sprintf("+book:%s +(%s)", bk1.Name, args[1])
	case "elastic":
		query = fmt.Sprintf("book:%s AND (%s)", bk1.Name, args[1])
	}

	_, total, err := idxConn.SearchNote(query, 1, 0)
	exitOnError(err)

	if total == 0 {
		fmt.Println("There were no notes that matched you query")
		return
	}

	cMsg := "This will move %d Notes from %s to %s, are you sure?"
	if skipConfirm || utils.AskForConfirmationMust(fmt.Sprintf(cMsg, total, bk1.Name, args[0])) {
		bk2, err := dbConn.GetOrCreateBookByName(args[0])
		exitOnError(err)

		var offset uint64
		for {
			ids, total, err := idxConn.SearchNote(query, 2048, 0)
			exitOnError(err)

			err = dbConn.EditNoteByIDBook(ids, bk2)
			exitOnError(err)

			notes, err := dbConn.GetAllNotesByIDs(ids)
			exitOnError(err)

			err = idxConn.IndexNotes(notes)
			exitOnError(err)

			offset = offset + uint64(len(ids))
			if offset >= total {
				break
			}
		}
	}
}

// SplitBookIDsCmd splits one book into two book using QueryStringQuery
var SplitBookIDsCmd = &cobra.Command{
	Use:   "ids [flags] <book_name> <note_id...>",
	Short: "Moves the given <note_id...> to a new or existing Book <book_name>",
	Long:  `Moves the Notes for the given <note_id...> into a new or existing Book.`,
	Run:   splitBooksIDsCmdRun,
}

func splitBooksIDsCmdRun(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		exitValidationError("invalid arguments given", cmd)
	}

	bk1 := workingNotebook

	ids := make([]int64, 0, len(args))
	for _, a := range args[1:] {
		noteID, err := strconv.ParseInt(a, 10, 64)
		exitOnError(err)
		ids = append(ids, noteID)
	}

	cMsg := "This will move %d Notes from %s to %s, are you sure?"
	if skipConfirm || utils.AskForConfirmationMust(fmt.Sprintf(cMsg, len(ids), bk1.Name, args[0])) {
		bk2, err := dbConn.GetOrCreateBookByName(args[0])
		exitOnError(err)

		err = dbConn.EditNoteByIDBook(ids, bk2)
		exitOnError(err)

		notes, err := dbConn.GetAllNotesByIDs(ids)
		exitOnError(err)

		err = idxConn.IndexNotes(notes)
		exitOnError(err)
	}
}
