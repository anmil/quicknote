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

package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/anmil/quicknote"

	"github.com/fatih/color"
)

// Terminal Colors
var (
	FgBlack   = color.New(color.FgBlack).SprintFunc()
	FgRed     = color.New(color.FgRed).SprintFunc()
	FgGreen   = color.New(color.FgGreen).SprintFunc()
	FgYellow  = color.New(color.FgYellow).SprintFunc()
	FgBlue    = color.New(color.FgBlue).SprintFunc()
	FgMagenta = color.New(color.FgMagenta).SprintFunc()
	FgCyan    = color.New(color.FgCyan).SprintFunc()
	FgWhite   = color.New(color.FgWhite).SprintFunc()

	FgHiBlack   = color.New(color.FgHiBlack).SprintFunc()
	FgHiRed     = color.New(color.FgHiRed).SprintFunc()
	FgHiGreen   = color.New(color.FgHiGreen).SprintFunc()
	FgHiYellow  = color.New(color.FgHiYellow).SprintFunc()
	FgHiBlue    = color.New(color.FgHiBlue).SprintFunc()
	FgHiMagenta = color.New(color.FgHiMagenta).SprintFunc()
	FgHiCyan    = color.New(color.FgHiCyan).SprintFunc()
	FgHiWhite   = color.New(color.FgHiWhite).SprintFunc()

	BgBlack   = color.New(color.BgBlack).SprintFunc()
	BgRed     = color.New(color.BgRed).SprintFunc()
	BgGreen   = color.New(color.BgGreen).SprintFunc()
	BgYellow  = color.New(color.BgYellow).SprintFunc()
	BgBlue    = color.New(color.BgBlue).SprintFunc()
	BgMagenta = color.New(color.BgMagenta).SprintFunc()
	BgCyan    = color.New(color.BgCyan).SprintFunc()
	BgWhite   = color.New(color.BgWhite).SprintFunc()

	BgHiBlack   = color.New(color.BgHiBlack).SprintFunc()
	BgHiRed     = color.New(color.BgHiRed).SprintFunc()
	BgHiGreen   = color.New(color.BgHiGreen).SprintFunc()
	BgHiYellow  = color.New(color.BgHiYellow).SprintFunc()
	BgHiBlue    = color.New(color.BgHiBlue).SprintFunc()
	BgHiMagenta = color.New(color.BgHiMagenta).SprintFunc()
	BgHiCyan    = color.New(color.BgHiCyan).SprintFunc()
	BgHiWhite   = color.New(color.BgHiWhite).SprintFunc()
)

// PrintNotes prints notes in the given format
func PrintNotes(notes quicknote.Notes, format string) error {
	var err error
	switch format {
	case "ids":
		PrintNotesIDs(notes)
	case "text":
		PrintNotesColored(notes, false)
	case "short":
		PrintNotesColored(notes, true)
	case "csv":
		err = PrintNotesCSV(notes)
	case "json":
		err = PrintNotesJSON(notes)
	}
	return err
}

// PrintNoteColored prints the Note to stdout in color
func PrintNoteColored(n *quicknote.Note, titleOnly bool) {
	if titleOnly {
		printNoteTitleOnly(n)
	} else {
		printDetailedNoteColored(n)
	}
}

// PrintNotesColored prints the list of Notes to stdout in color
func PrintNotesColored(notes quicknote.Notes, titleOnly bool) {
	if titleOnly {
		printNotesTitleOnly(notes)
	} else {
		printDetailedNotes(notes)
	}
}

func printNotesTitleOnly(notes quicknote.Notes) {
	for _, n := range notes {
		printNoteTitleOnly(n)
	}
}

func printNoteTitleOnly(n *quicknote.Note) {
	fmt.Print(FgCyan("ID: "))
	fmt.Print(FgMagenta(n.ID))
	fmt.Print(FgCyan(" Title: "))
	fmt.Println(n.Title)
}

func printDetailedNotes(notes quicknote.Notes) {
	nLen := len(notes)
	for idx, n := range notes {
		printDetailedNoteColored(n)
		if idx+1 < nLen {
			fmt.Printf("\n")
		}
	}
}

func printDetailedNoteColored(n *quicknote.Note) {
	fmt.Print(FgCyan("ID: "))
	fmt.Print(n.ID)

	fmt.Print(FgCyan(" Book: "))
	fmt.Print(n.Book.Name)

	fmt.Print(FgCyan(" Type: "))
	fmt.Print(n.Type)

	fmt.Print(FgCyan(" Created: "))
	fmt.Print(n.Created.Format("2006-01-02 03:04:05 PM"))

	fmt.Print(FgCyan(" Modified: "))
	fmt.Println(n.Modified.Format("2006-01-02 03:04:05 PM"))

	fmt.Println(FgMagenta("--------------------------------------------------"))

	fmt.Print(FgCyan("Title: "))
	fmt.Println(n.Title)

	fmt.Print(FgCyan("Tags: "))
	fmt.Println(strings.Join(colorTags(n.Tags), ", "))

	if len(n.Body) > 0 {
		fmt.Printf("\n%s\n", n.Body)
	}
}

func colorTags(tags quicknote.Tags) []string {
	ctags := make([]string, 0, len(tags))
	for i := 0; i < len(tags); i++ {
		ctags = append(ctags, FgBlue(tags[i].Name))
	}
	return ctags
}

// PrintNotesIDs prints the Note's ids
func PrintNotesIDs(notes quicknote.Notes) {
	for _, n := range notes {
		fmt.Println(n.ID)
	}
}

// PrintNotesCSV prints Notes in csv format
func PrintNotesCSV(notes quicknote.Notes) error {
	w := csv.NewWriter(os.Stdout)
	w.Write([]string{"id", "created", "modified", "type", "title", "body", "book", "tags"})

	for _, n := range notes {
		err := w.Write([]string{
			strconv.FormatInt(n.ID, 10),
			n.Created.Format("2006-01-02 03:04:05 PM"),
			n.Modified.Format("2006-01-02 03:04:05 PM"),
			n.Type,
			n.Title,
			n.Body,
			n.Book.Name,
			strings.Join(n.GetTagStringArray(), ", "),
		})
		if err != nil {
			return err
		}
	}

	w.Flush()
	err := w.Error()
	return err
}

// PrintNotesJSON prints Notes in json format
func PrintNotesJSON(notes quicknote.Notes) error {
	b, err := json.Marshal(notes)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
