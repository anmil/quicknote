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
	"time"

	"github.com/anmil/quicknote/note"
	"github.com/spf13/cobra"
)

func init() {
	NewCmd.AddCommand(NewbookCmd)
}

// NewbookCmd Create a new Book
var NewbookCmd = &cobra.Command{
	Use:   "book <book_name...>",
	Short: "Create a new Book",
	Long: `Create a new Book

Books allow you to organize collections of notes. Every Note must belong to a
Book. All commends (unless stated otherwise) operates only on one book. Such as,
if you call 'qnote ls notes', it only list the notes for the working Book. The
working Book can be changed with the '-n' flag, or you can changed the default
Book in the config file. In most cases, it is not advised to create to many
books, but they are useful to keeping work notes separated from personal.`,
	Run: newbookCmdRun,
}

func newbookCmdRun(cmd *cobra.Command, args []string) {
	for _, name := range args {
		bk, err := dbConn.GetBookByName(name)
		exitOnError(err)
		if bk != nil {
			fmt.Printf("Notebook %s already existed\n", bk.Name)
		} else {
			bk = &note.Book{
				Created:  time.Now(),
				Modified: time.Now(),
				Name:     name,
			}

			err = dbConn.CreateBook(bk)
			exitOnError(err)

			fmt.Printf("Notebook %s created\n", bk.Name)
		}
	}
}
