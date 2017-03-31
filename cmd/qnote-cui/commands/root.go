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
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/anmil/quicknote"
	"github.com/anmil/quicknote/cmd/shared/config"
	"github.com/jroimartin/gocui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Version information
var (
	VersionMajor    = 0
	VersionMinor    = 1
	VersionRevision = 0
)

var (
	dbConn          quicknote.DB
	idxConn         quicknote.Index
	workingNotebook *quicknote.Book
)

var (
	workingNotebookName string
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&workingNotebookName, "notebook", "n",
		viper.GetString("default_notebook"), "Working Notebook")
}

// RootCmd is the base command
var RootCmd = &cobra.Command{
	Use:               "qnote-cui",
	Short:             "",
	PersistentPreRun:  PreseistentPreRunRoot,
	PersistentPostRun: PreseistentPostRunRoot,
	Run:               rootCmdRun,
}

func rootCmdRun(cmd *cobra.Command, args []string) {
	g, err := gocui.NewGui(gocui.Output256)

	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.InputEsc = true

	g.SetManagerFunc(mainLayout)

	setKeybindings(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

// PreseistentPreRunRoot runs before the Root Command and any child
// commands that do not override it.
func PreseistentPreRunRoot(cmd *cobra.Command, args []string) {
	var err error
	dbConn, err = config.GetDBConn()
	exitOnError(err)
	idxConn, err = config.GetIndexConn()
	exitOnError(err)

	workingNotebook, err = config.GetWorkingBook(dbConn, workingNotebookName)
	exitOnError(err)

	if workingNotebook == nil {
		exitOnError(errors.New("Notebook does not exists"))
	}
}

// PreseistentPostRunRoot runs after the Root Command and any child
// commands that do not override it.
func PreseistentPostRunRoot(cmd *cobra.Command, args []string) {
	if dbConn != nil {
		dbConn.Close()
	}
}

func exitOnError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func exitValidationError(msg string, cmd *cobra.Command) {
	fmt.Println(msg)
	cmd.Usage()
	os.Exit(1)
}

func setKeybindings(g *gocui.Gui) {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quitCB); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("search_box", gocui.KeyEsc, gocui.ModNone, quitCB); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("search_box", gocui.KeyEnter, gocui.ModNone, displayNote); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("search_box", gocui.KeyArrowUp, gocui.ModNone, searchBoxKeyUpEvent); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("search_box", gocui.KeyArrowDown, gocui.ModNone, searchBoxKeyDownEvent); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("note_display", gocui.KeyEsc, gocui.ModNone, delDisplayNote); err != nil {
		log.Panicln(err)
	}
}

func quitCB(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

var defaultConfigFileText = ``
