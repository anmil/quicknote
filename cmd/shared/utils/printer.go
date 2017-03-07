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
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/anmil/quicknote/note"
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
func PrintNotes(notes []*note.Note, format string) error {
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
func PrintNoteColored(notes *note.Note, titleOnly bool) {
	if titleOnly {
		printNoteTitleOnly(notes)
	} else {
		printDetailedNoteColored(notes)
	}
}

// PrintNotesColored prints the list of Notes to stdout in color
func PrintNotesColored(notes []*note.Note, titleOnly bool) {
	if titleOnly {
		printNotesTitleOnly(notes)
	} else {
		printDetailedNotes(notes)
	}
}

func printNotesTitleOnly(notes []*note.Note) {
	for _, n := range notes {
		printNoteTitleOnly(n)
	}
}

func printNoteTitleOnly(n *note.Note) {
	fmt.Print(FgCyan("ID: "))
	fmt.Print(FgMagenta(n.ID))
	fmt.Print(FgCyan(" Title: "))
	fmt.Println(n.Title)
}

func printDetailedNotes(notes []*note.Note) {
	nLen := len(notes)
	for idx, n := range notes {
		printDetailedNoteColored(n)
		if idx+1 < nLen {
			fmt.Printf("\n")
		}
	}
}

func printDetailedNoteColored(n *note.Note) {
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

func colorTags(tags []*note.Tag) []string {
	ctags := make([]string, 0, len(tags))
	for i := 0; i < len(tags); i++ {
		ctags = append(ctags, FgBlue(tags[i].Name))
	}
	return ctags
}

// PrintNotesIDs prints the Note's ids
func PrintNotesIDs(notes []*note.Note) {
	for _, n := range notes {
		fmt.Println(n.ID)
	}
}

// PrintNotesCSV prints Notes in csv format
func PrintNotesCSV(notes []*note.Note) error {
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
func PrintNotesJSON(notes []*note.Note) error {
	b, err := json.Marshal(notes)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

// longestColumnWordString simple loop to find the longest word in the slice
func longestColumnWordString(ws []string) int {
	lw := 0
	for _, w := range ws {
		if len(w) > lw {
			// l := len(fmt.Sprint(w))
			l := utf8.RuneCountInString(w)
			lw = l
		}
	}
	return lw
}

// buildGridString builds the grid from the slice of words. The grid
// is sorted from top down than left to right. The number of columns
// is calculated using the len of the slice and number of rows given
func buildGridString(words []string, rCnt int, cb func(string) string) ([][]string, string, int) {

	// Calculate the number of columns needed to
	// get have the corrected row count
	d := float64(len(words)) / float64(rCnt)
	cCnt := int(math.Ceil(d))

	// Using the rCnt and cCnt, we build the
	// two dimensional table to represent the grid
	table := make([][]string, 0)
	for c := 0; c < cCnt; c++ {
		col := make([]string, 0)
		sIdx := c * rCnt
		if sIdx < len(words) {
			eIdx := sIdx + rCnt
			for r := sIdx; r < eIdx; r++ {
				if r >= len(words) {
					break
				}
				col = append(col, words[r])
			}
			table = append(table, col)
		}
	}

	// To keep all of the column aligned, we find the
	// longest word in each column to use as our padding
	// value
	paddingColMap := make(map[int]int)
	for idx, col := range table {
		paddingColMap[idx] = longestColumnWordString(col)
	}

	// Now we make the grid by printing horizontally across
	// the table. If the word is smaller than the longest
	// word in the column (using paddingColMap). Than we
	// pad the word with spaces till its of equal length.
	// This ensures that all of the columns are aligned.
	// We also keep track of the longest rows (in character len)
	// and return it and the string containing the grid.
	mRlen := 0
	var msg string
	for r := 0; r < rCnt; r++ {
		line := ""
		for c := 0; c < len(table); c++ {
			if r >= len(table[c]) {
				break
			}
			lw := paddingColMap[c]
			w := cb(table[c][r])
			wLen := len(table[c][r])
			pw := fmt.Sprintf("%s%s ", w, strings.Repeat(" ", lw-wLen))
			line = fmt.Sprintf("%s%s", line, pw)
		}
		if len(line) > mRlen {
			mRlen = len(line)
		}
		msg = fmt.Sprintf("%s%s\n", msg, line)
	}

	return table, msg, mRlen
}

// see BuildGridString for details
func buildGridStringRec(strs []string, start, end, maxLen int, m string, table [][]string, cb func(string) string) ([][]string, string) {
	if end-start <= 1 {
		return table, m
	}

	p := ((end - start) / 2) + start
	table, msg, mRlen := buildGridString(strs, p, cb)

	if mRlen < maxLen {
		return buildGridStringRec(strs, start, p, maxLen, msg, table, cb)
	}
	return buildGridStringRec(strs, p, end, maxLen, m, table, cb)
}

// BuildGridString calls a binary search style recursive function that finds
// the minimum number of rows without passing the maximum column width.
// Generally, maxLen will be equal to the width of the terminal (or view if using
// a CUI).
func BuildGridString(strs []string, maxLen int) ([][]string, string) {
	sort.Strings(strs)
	cb := func(w string) string {
		return w
	}
	return buildGridStringRec(strs, 0, len(strs), maxLen, "", nil, cb)
}

// BuildGridStringCB like BuildGridString but accepts a call back function to alter
// the word before adding it to the final grid. Useful to adding things like color
// text that would normally message column lengths.
func BuildGridStringCB(strs []string, maxLen int, cb func(string) string) ([][]string, string) {
	sort.Strings(strs)
	return buildGridStringRec(strs, 0, len(strs), maxLen, "", nil, cb)
}
