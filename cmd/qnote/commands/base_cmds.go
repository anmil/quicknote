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
	"strings"

	"github.com/anmil/quicknote/cmd/shared/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Command line variables
var (
	workingNotebookName string
	displayOrder        string
	sortBy              string
	displayFormat       string
	skipConfirm         bool
)

var displayOrderOptions = []string{
	"asc",
	"desc",
}

var sortByOptinos = []string{
	"id",
	"created",
	"modified",
	"title",
}

var displayFormatOptions = []string{
	"ids",
	"text",
	"short",
	"json",
	"csv",
}

func init() {
	RootCmd.AddCommand(NewCmd)
	RootCmd.AddCommand(GetCmd)
	RootCmd.AddCommand(EditCmd)
	RootCmd.AddCommand(DeleteCmd)
	RootCmd.AddCommand(SplitCmd)

	viper.SetDefault("display_order", "asc")
	viper.SetDefault("order_by", "modified")
	viper.SetDefault("display_format", "text")

	GetCmd.PersistentFlags().StringVarP(&displayOrder, "display-order", "d", viper.GetString("display_order"),
		fmt.Sprintf("The order to display Notebook, Notes, Tags [%s]", strings.Join(displayOrderOptions, ", ")))

	GetCmd.PersistentFlags().StringVarP(&sortBy, "sort-by", "s", viper.GetString("order_by"),
		fmt.Sprintf("Sort notes by [%s]", strings.Join(sortByOptinos, ", ")))

	GetCmd.PersistentFlags().StringVarP(&displayFormat, "display-format", "f", viper.GetString("display_format"),
		fmt.Sprintf("Format to display notes in [%s]", strings.Join(displayFormatOptions, ", ")))

	DeleteCmd.PersistentFlags().BoolVarP(&skipConfirm, "skip-confirm", "", false, "Do not prompt to confirm action")
}

// NewCmd create new Note or Notebook
var NewCmd = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create", "add"},
	Short:   "Create new Note or Notebook",
}

// GetCmd get/list Notes, Notebook, Tags
var GetCmd = &cobra.Command{
	Use:              "get",
	Aliases:          []string{"list", "ls"},
	Short:            "Get/List Notes, Notebook, Tags",
	PersistentPreRun: preseistentPreGetRoot,
}

func preseistentPreGetRoot(cmd *cobra.Command, args []string) {
	PreseistentPreRunRoot(cmd, args)
	validateGetFlags(cmd)
}

func validateGetFlags(cmd *cobra.Command) {
	if !utils.InSliceString(displayOrder, displayOrderOptions) {
		exitValidationError("invalid display-order", cmd)
	}
	if !utils.InSliceString(sortBy, sortByOptinos) {
		exitValidationError("invalid display-order", cmd)
	}
	if !utils.InSliceString(displayFormat, displayFormatOptions) {
		exitValidationError("invalid display-format", cmd)
	}
}

// EditCmd edit Note or Notebook
var EditCmd = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"update"},
	Short:   "Edit Note or Notebook",
}

// DeleteCmd delete Note or Notebook
var DeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"del", "remove", "rm"},
	Short:   "Delete Note or Notebook",
}

// SplitCmd splits Books
var SplitCmd = &cobra.Command{
	Use:   "split",
	Short: "Split Book",
}
