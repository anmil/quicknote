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
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/anmil/quicknote/cmd/shared/encoding"
	"github.com/anmil/quicknote/cmd/shared/utils"

	"github.com/spf13/cobra"
)

var (
	compressOutput bool
	outputFile     string
)

func init() {
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

	var out io.Writer
	out = os.Stdout

	fp, fn, err := getFilePath(outputFile, compressOutput)
	exitOnError(err)

	if len(outputFile) > 0 {
		file, err := os.Create(fp)
		exitOnError(err)
		defer file.Close()

		out = file
	}

	if compressOutput {
		zip := gzip.NewWriter(out)
		zip.Name = fn
		zip.Comment = "Exported notes form QuickNote Qnote"
		zip.ModTime = time.Now()
		defer zip.Close()

		out = zip
	}

	enc := encoding.NewBinaryEncoder(out)
	_, err = enc.WriteHeader()
	exitOnError(err)

	for _, n := range notes {
		_, err = enc.WriteNote(n)
		exitOnError(err)
	}
}

func getFilePath(outputFile string, compressed bool) (string, string, error) {
	fp := "notes.qnot"
	fn := fp

	if len(outputFile) > 0 {
		fp, err := utils.ExpandFilePath(outputFile)
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
