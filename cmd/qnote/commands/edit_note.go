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
	"time"

	"github.com/anmil/quicknote/cmd/shared/utils"
	"github.com/anmil/quicknote/note"
	"github.com/anmil/quicknote/parser"
	"github.com/spf13/cobra"
)

func init() {
	EditCmd.AddCommand(EditNoteCmd)
	EditNoteCmd.AddCommand(MoveNotesIDsCmd)
}

// EditNoteCmd Edit note
var EditNoteCmd = &cobra.Command{
	Use:   "note <note id>",
	Short: "Edit note",
	Long:  `Opens an editor (default vim) to allow you to edit a Note`,
	Run:   editNoteCmdRun,
}

func editNoteCmdRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		exitValidationError("No Note ID given", cmd)
	}

	noteID, err := strconv.ParseInt(args[0], 10, 64)
	exitOnError(err)

	oldNote, err := dbConn.GetNoteByID(noteID)
	exitOnError(err)

	if oldNote == nil {
		fmt.Println("Note does not exists")
		return
	}

	editor, err := utils.NewEditor()
	exitOnError(err)
	defer editor.Close()

	editor.SetText(fmt.Sprintf("%s\n%s", oldNote.Title, oldNote.Body))
	err = editor.Open()
	exitOnError(err)

	p, err := parser.NewParser(oldNote.Type)
	exitOnError(err)
	p.Parse(editor.Text())

	tags := make([]*note.Tag, 0, len(p.Tags()))
	for _, t := range p.Tags() {
		tag, err := dbConn.GetOrCreateTagByName(t)
		exitOnError(err)
		tags = append(tags, tag)
	}

	newNote := &note.Note{
		ID:       oldNote.ID,
		Created:  oldNote.Created,
		Modified: time.Now(),
		Book:     oldNote.Book,
		Type:     oldNote.Type,
		Title:    p.Title(),
		Body:     p.Body(),
		Tags:     tags,
	}

	err = dbConn.EditNote(newNote)
	exitOnError(err)

	err = idxConn.IndexNote(newNote)
	exitOnError(err)

	utils.PrintNoteColored(newNote, false)
}

// MoveNotesIDsCmd See SplitBookIDsCmd
var MoveNotesIDsCmd = &cobra.Command{
	Use:   "move [flags] <book_name> <note_id...>",
	Short: SplitBookIDsCmd.Short,
	Long:  SplitBookIDsCmd.Long,
	Run:   SplitBookIDsCmd.Run,
}
