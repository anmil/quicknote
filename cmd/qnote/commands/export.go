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
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/anmil/quicknote/cmd/shared/encoding"
	"github.com/anmil/quicknote/cmd/shared/utils"
	"github.com/anmil/quicknote/note"

	"github.com/spf13/cobra"
)

var (
	compressOutput bool
	outputFile     string
)

func init() {
	ExportCmd.AddCommand(ExportBookCmd)

	ExportCmd.PersistentFlags().StringVarP(&outputFile, "out-file", "o", "", "Write to file instead of stdout")
	ExportCmd.PersistentFlags().BoolVarP(&compressOutput, "compress", "c", false, "Compress output with gzip")
}

// ExportCmd Export all Notes, Books, and Tags
var ExportCmd = &cobra.Command{
	Use:   "export [flags]",
	Short: "Export all Notes, Books, and Tags",
	Long: `Export all Notes, Books, and Tags using the QNOT file format.

See the documentation in github.com/anmil/quicknote/cmd/shard/encoding/binary.go
for the specifications of the format. This command is useful for backing up all
of your notes or transferring them to another system.

See the help docs for the "import" command for how importing is done and conflicts
are dealt with.

The output is written to stdout by default, use the '-o' flag to write to a file.

Qnote can also compress the output with gzip by specifying the '-c' flag.`,
	Run: exportCmdRun,
}

func exportCmdRun(cmd *cobra.Command, args []string) {
	notes, err := dbConn.GetAllNotes("created", "asc")
	exitOnError(err)

	out, fn, file, err := getExportWriter()
	exitOnError(err)
	if file != nil {
		defer file.Close()
	}

	err = exportNotes(notes, fn, out, compressOutput)
	exitOnError(err)
}

// ExportBookCmd Export all Notes, Tags in book(s)
var ExportBookCmd = &cobra.Command{
	Use:   "book [flags] <book>...",
	Short: "Export all Notes, Tags in book(s)",
	Long: `Export all Notes, Tags in book(s) using the QNOT file format.

See the help documentation for the export command for details

	qnote help export

`,
	Run: exportBookCmdRun,
}

func exportBookCmdRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		exitValidationError("No books given", cmd)
	}

	var books note.Books
	for _, bkName := range args {
		bk, err := dbConn.GetBookByName(bkName)
		exitOnError(err)
		if bk == nil {
			fmt.Printf("Book %s does not exists\n", bkName)
			return
		}

		books = append(books, bk)
	}

	out, fn, file, err := getExportWriter()
	exitOnError(err)
	if file != nil {
		defer file.Close()
	}

	var notes note.Notes
	for _, bk := range books {
		ns, err := dbConn.GetAllBookNotes(bk, "created", "asc")
		exitOnError(err)
		notes = append(notes, ns...)
	}

	err = exportNotes(notes, fn, out, compressOutput)
	exitOnError(err)
}

func exportNotes(notes note.Notes, fileName string, out io.Writer, compressed bool) error {
	if compressed {
		zip := gzip.NewWriter(out)
		zip.Name = fileName
		zip.Comment = "Exported notes form QuickNote Qnote"
		zip.ModTime = time.Now()
		defer zip.Close()

		out = zip
	}

	enc := encoding.NewBinaryEncoder(out)
	if _, err := enc.WriteHeader(); err != nil {
		return err
	}

	for _, n := range notes {
		if _, err := enc.WriteNote(n); err != nil {
			return err
		}
	}

	return nil
}

func getExportWriter() (io.Writer, string, *os.File, error) {
	var out io.Writer
	out = os.Stdout

	fp, fn, err := getFilePath(outputFile, compressOutput)
	if err != nil {
		return nil, "", nil, err
	}

	var file *os.File
	if len(outputFile) > 0 {
		file, err = os.Create(fp)
		if err != nil {
			return nil, "", nil, err
		}
		out = file
	}

	return out, fn, file, nil
}

func getFilePath(outputFile string, compressed bool) (string, string, error) {
	fp := "notes.qnot"
	fn := fp

	var err error
	if len(outputFile) > 0 {
		fp, err = utils.ExpandFilePath(outputFile)
		if err != nil {
			return "", "", err
		}
		_, fn = filepath.Split(fp)
	}

	if compressed {
		if utils.GetFileExt(fp) != ".gz" {
			fp = fp + ".gz"
		}
		if utils.GetFileExt(fn) == ".gz" {
			fn = string(fn[:len(fn)-3])
		}
	}

	return fp, fn, nil
}
