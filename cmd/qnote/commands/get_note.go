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
	"strconv"

	"github.com/anmil/quicknote/cmd/shared/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	GetCmd.AddCommand(GetNoteCmd)
	GetNoteCmd.AddCommand(GetNoteAllCmd)

	viper.SetDefault("titles_only", "false")
}

// GetNoteCmd Gets all Notes, or Notes for the given IDs
var GetNoteCmd = &cobra.Command{
	Use:     "note [flags] [note id...]",
	Aliases: []string{"notes"},
	Short:   "List all notes in the working Book, or all notes for the given [note id...]",
	Long: `List all Notes for the Book, or Notes for the given IDs.

Prints Notes in the format given by '-f', see '-f' docs for all available
options.`,
	Run: getNoteCmdRun,
}

func getNoteCmdRun(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		noteIDs := make([]int64, 0, len(args))
		for _, a := range args {
			noteID, err := strconv.ParseInt(a, 10, 64)
			exitOnError(err)
			noteIDs = append(noteIDs, noteID)
		}

		notes, err := dbConn.GetAllNotesByIDs(noteIDs)
		exitOnError(err)

		err = utils.PrintNotes(notes, displayFormat)
		exitOnError(err)
	} else {
		getAllBookNotes()
	}
}

func getAllBookNotes() {
	notes, err := dbConn.GetAllBookNotes(workingNotebook, sortBy, displayOrder)
	exitOnError(err)
	err = utils.PrintNotes(notes, displayFormat)
	exitOnError(err)
}

// GetNoteAllCmd Gets all Notes
var GetNoteAllCmd = &cobra.Command{
	Use:   "all",
	Short: "List all Notes for all Books",
	Long: `List all notes in all Books

This is the same as 'gnote ls notes' except it returns all Notes in all Books`,
	Run: getNoteAllCmdRun,
}

func getNoteAllCmdRun(cmd *cobra.Command, args []string) {
	notes, err := dbConn.GetAllNotes(sortBy, displayOrder)
	exitOnError(err)
	err = utils.PrintNotes(notes, displayFormat)
	exitOnError(err)
}
