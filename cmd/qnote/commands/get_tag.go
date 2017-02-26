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

	"github.com/spf13/cobra"
)

func init() {
	GetCmd.AddCommand(GetTagCmd)
}

// GetTagCmd lists all tags for the working Book
var GetTagCmd = &cobra.Command{
	Use:     "tag",
	Aliases: []string{"tags"},
	Short:   "lists all tags for the working Book",
	Run:     getTagCmdRun,
}

func getTagCmdRun(cmd *cobra.Command, args []string) {
	tags, err := dbConn.GetAllBookTags(workingNotebook)
	exitOnError(err)

	for _, tag := range tags {
		fmt.Println(tag.Name)
	}
}
