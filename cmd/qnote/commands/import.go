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
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/anmil/quicknote"
	"github.com/anmil/quicknote/cmd/shared/encoding"
	"github.com/anmil/quicknote/cmd/shared/utils"

	"github.com/spf13/cobra"
)

var (
	skipDupCheck     bool
	preserveModified bool
	inputCompressed  bool
)

func init() {
	ImportCmd.PersistentFlags().BoolVarP(&skipDupCheck, "skip-dup-check", "s", false,
		"Skip duplicate Notes (slows import down)")
	ImportCmd.PersistentFlags().BoolVarP(&preserveModified, "preserve-modified", "p", false,
		"If set, the note's modified from the QNOT file is not overridden")
	ImportCmd.PersistentFlags().BoolVarP(&inputCompressed, "compress", "c", false, "Input from stdin or file is gzip compressed")
}

// ImportCmd Imports Notes, Books, and Tags from a QNOT file
var ImportCmd = &cobra.Command{
	Use:   "import [flags] [<input-file>]",
	Short: "Imports Notes, Books, and Tags from a QNOT file",
	Long: `Imports Notes, Books, and Tags from a QNOT file

Imports notes from a QNOT file that is read from stdin or a file passed
as an argument.

If the QNOT file contains a book that is already in the database it is not
recreated and the existing book is used. The same goes for Tags.

By default the importer checks if a Note already exists. A Note is considered
to be equal if the Book, Type, Title, and Body are the same. If a duplicate
is found the Note is skipped. The duplicate check can be disabled with the
"--skip-dup-check" flag in which case all Notes are saved as new Notes.

Skipping the check is a good idea if you know there are no duplicates -- the
importer will run faster since it does not have to search the database for
existing Notes.

A Note's ID is not preserved
Created dates are preserved
Modified dates are set to the current time (set --preserve-modified to disable this)

Be sure to set the "-c" flag if the QNOT file is compressed with gzip. If you
pass the file name in as an argument and it ends with ".gz" Qnote will automatically
treat it as compressed.`,
	Run: importCmdRun,
}

func importCmdRun(cmd *cobra.Command, args []string) {
	var in io.Reader
	in = os.Stdin

	if len(args) > 0 {
		file, err := os.Open(args[0])
		exitOnError(err)
		defer file.Close()
		in = file

		// Auto detect if input file is compressed
		ext := utils.GetFileExt(file.Name())
		if ext == ".gz" {
			inputCompressed = true
		}
	}

	if inputCompressed {
		zip, err := gzip.NewReader(in)
		exitOnError(err)
		defer zip.Close()
		in = zip
	}

	r := bufio.NewReader(in)
	dec := encoding.NewBinaryDecoder(r)
	err := dec.ParseHeader()
	exitOnError(err)

	fmt.Println("Version:", dec.Header.Version)
	fmt.Println("Created:", dec.Header.Created)

	notes, err := dec.ParseNotes()
	exitOnError(err)

	bkNew := make(map[string]bool)
	books := make(map[string]*quicknote.Book)
	tags := make(map[string]*quicknote.Tag)

	for n := range notes {
		bk, found := books[n.Book.Name]
		if !found {
			bk, err = dbConn.GetBookByName(n.Book.Name)
			exitOnError(err)

			if bk == nil {
				err = dbConn.CreateBook(n.Book)
				exitOnError(err)
				bk = n.Book

				books[bk.Name] = bk
				bkNew[bk.Name] = true
			} else {
				bkNew[bk.Name] = false
			}
		}
		n.Book = bk

		// When checking for duplicate notes, if the book is new.
		// We know this can not be a duplicate.
		if isNew, _ := bkNew[bk.Name]; !isNew && !skipDupCheck {
			// Save this values as the lookup will destroys them if
			// the note does not exists
			n.ID = -1
			created := n.Created
			err = dbConn.GetNoteByNote(n)
			exitOnError(err)

			if n.ID > 0 {
				fmt.Print("Skipping Dup: ")
				utils.PrintNoteColored(n, true)
				continue
			}

			if n.Created != created {
				fmt.Println("Lost created date")
			}
		}

		for i := 0; i < len(n.Tags); i++ {
			tag, found := tags[n.Tags[i].Name]
			if !found {
				tag, err = dbConn.GetTagByName(n.Tags[i].Name)
				exitOnError(err)

				if tag == nil {
					err = dbConn.CreateTag(n.Tags[i])
					exitOnError(err)
					tags[tag.Name] = tag
				}
			}
			n.Tags[i] = tag
		}

		if !preserveModified {
			n.Modified = time.Now()
		}

		err = saveNote(n)
		exitOnError(err)

		fmt.Print("Saved Note: ")
		utils.PrintNoteColored(n, true)
	}
	exitOnError(dec.Err)
}
