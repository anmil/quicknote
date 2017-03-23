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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/anmil/quicknote/cmd/shared/utils"
	"github.com/anmil/quicknote/note"
	"github.com/anmil/quicknote/parser"
	"github.com/spf13/cobra"
)

var (
	noteType string
)

var editorDocMessage = `%s

# Please enter the text for the note you wish to create. Empty notes
# aborts the creation. All lines below this message are ignored.
#
# First line is used as the title. Any word that starts with '#' is
# considered a tag
#
# This note will be saved with the following values:
#      Notebook: %s
#          Type: %s
#`

func init() {
	NewCmd.AddCommand(NewNoteCmd)
	NewCmd.AddCommand(NewURLNoteCmd)
	NewCmd.AddCommand(NewNoteFromJSONCmd)

	NewNoteCmd.Flags().StringVarP(&noteType, "note-type", "t", "basic",
		fmt.Sprintf("The new Note's type [%s]", strings.Join(note.NoteTypes, ", ")))
}

// NewNoteCmd Create a new basic note
var NewNoteCmd = &cobra.Command{
	Use:   "note",
	Short: "Create a new basic Note",
	Long: `Create a new note to store and index

Opens an editor (default vim) to allow you to enter a new Note. The Note text
is parsed using the first line as the Note's title. All other lines are used
as the Note's body. Any word starting with '#' character is parsed as a Tag.
Tags can be in either the title or the body.`,
	Run: newNoteCmdRun,
}

func newNoteCmdRun(cmd *cobra.Command, args []string) {
	validateNewNoteFlags(cmd)

	switch noteType {
	case note.URL:
		newURLNoteCmdRun(cmd, args)
	default:
		editorText := fmt.Sprintf(editorDocMessage, "", workingNotebook.Name, note.Basic)
		createNewNote(editorText, note.Basic)
	}
}

func validateNewNoteFlags(cmd *cobra.Command) {
	if !utils.InSliceString(noteType, note.NoteTypes) {
		exitValidationError("invalid note type", cmd)
	}
}

// NewURLNoteCmd Create a new url note
var NewURLNoteCmd = &cobra.Command{
	Use:   "url <url>",
	Short: "Create a new url note",
	Long: `Creates a new note with URL type

DO NOT USE THIS IF YOU DO NOT WANT A REQUEST SENT TO THE URL

This is the same as calling 'qnote note new -t url <url>

This command preforms a GET request on the given URL. If a valid response
is given, it will parse the HTML for the 'title', keywords', and 'description'
meta tags and pre-fills the editor with this information.
`,
	Run: newURLNoteCmdRun,
}

func newURLNoteCmdRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		exitValidationError("invalid arguments", cmd)
	}

	url := args[0]
	text := getURLMetaNote(url)
	editorText := fmt.Sprintf(editorDocMessage, text, workingNotebook.Name, note.URL)
	createNewNote(editorText, note.URL)
}

func getURLMetaNote(url string) string {
	doc, err := goquery.NewDocument(url)
	exitOnError(err)

	var text string
	var keywords string
	var description string

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, found := s.Attr("name"); !found {
			return
		} else if name == "keywords" {
			keywords = getMetaContent(s)
		} else if name == "description" {
			description = getMetaContent(s)
		}
	})

	title := doc.Find("title").Text()
	if len(title) > 0 {
		text = fmt.Sprintf("%s - %s", title, url)
	} else {
		text = url
	}

	tags := parseTagsFromKeywords(keywords)

	if len(tags) > 0 {
		text = fmt.Sprintf("%s\n%s", text, strings.Join(tags, ", "))
	}
	if len(description) > 0 {
		text = fmt.Sprintf("%s\n\n%s", text, description)
	}

	return text
}

func getMetaContent(s *goquery.Selection) string {
	if content, found := s.Attr("content"); found {
		return content
	}
	return ""
}

func parseTagsFromKeywords(keywords string) []string {
	words := strings.Split(keywords, ",")
	tags := make([]string, 0, len(words))

	for _, w := range words {
		w = strings.TrimSpace(w)
		if len(w) > 0 {
			w = strings.Replace(w, " ", "_", -1)
			tags = append(tags, fmt.Sprintf("#%s", w))
		}
	}

	return tags
}

func createNewNote(text string, typ string) {
	editor, err := utils.NewEditor()
	exitOnError(err)
	defer editor.Close()

	if len(text) > 0 {
		editor.SetText(text)
	}

	err = editor.Open()
	exitOnError(err)

	noteText := editor.Text()
	noteText = removeBottomComment(noteText)
	noteText = strings.TrimSpace(noteText)
	if len(noteText) == 0 {
		fmt.Println("No text entered.. aborting")
		return
	}

	p, err := parser.NewParser(note.Basic)
	exitOnError(err)
	p.Parse(noteText)

	tags := make(note.Tags, 0, len(p.Tags()))
	for _, t := range p.Tags() {
		tag, err := dbConn.GetOrCreateTagByName(t)
		exitOnError(err)
		tags = append(tags, tag)
	}

	n := &note.Note{
		Created:  time.Now(),
		Modified: time.Now(),
		Book:     workingNotebook,
		Type:     typ,
		Title:    p.Title(),
		Body:     p.Body(),
		Tags:     tags,
	}

	err = saveNote(n)
	exitOnError(err)

	utils.PrintNoteColored(n, false)
}

func removeBottomComment(text string) string {
	text = strings.TrimRight(text, "\n")
	lines := strings.Split(text, "\n")

	endIndex := len(lines) - 1
	for ; endIndex >= 0; endIndex-- {
		if !strings.HasPrefix(lines[endIndex], "#") {
			break
		}
	}

	// If there is nothing above the button
	// comments, than endIdex ends up -1
	if endIndex < 0 {
		return ""
	}
	return strings.Join(lines[:endIndex+1], "\n")
}

// NewURLNoteCmd Create new notes from JSON
var NewNoteFromJSONCmd = &cobra.Command{
	Use:   "json [<json>]",
	Short: "Create new notes from JSON",
	Long: fmt.Sprintf(`Create new notes from JSON

JSON must be in the format

[
	...
	{
		"title": "<title>",
		"type": "<type>",
		"tags": ["<tag1>", "<tag2>", ...],
		"body": "<body>",
		"book": "<book>"
	},
	...
]

"type" must be one of the following: %s
If the book does not exists, it will be created

If <json> is note given, qnote will read from stdin
`, strings.Join(note.NoteTypes, ", ")),
	Run: newNoteFromJSONCmdRun,
}

type jNote struct {
	Title string   `json:"title"`
	Type  string   `json:"type"`
	Tags  []string `json:"tags"`
	Body  string   `json:"body"`
	Book  string   `json:"book"`
}

func newNoteFromJSONCmdRun(cmd *cobra.Command, args []string) {
	var reader io.Reader
	if len(args) >= 1 {
		reader = strings.NewReader(args[0])
	} else {
		reader = os.Stdin
	}

	var wg sync.WaitGroup
	var nCnt int64

	jnChan := make(chan *jNote, 1024)
	results := make(chan error, 2048)

	go func() {
		for r := range results {
			nCnt++
			if r != nil {
				fmt.Println("Error:", r)
			}
		}
	}()

	for i := 0; i < 16; i++ {
		wg.Add(1)
		go createJNoteWorker(i, &wg, jnChan, results)
	}

	dec := json.NewDecoder(reader)
	_, err := dec.Token()
	exitOnError(err)

	for dec.More() {
		var jn jNote
		err := dec.Decode(&jn)
		exitOnError(err)

		jnChan <- &jn
	}

	_, err = dec.Token()
	exitOnError(err)

	close(jnChan)
	wg.Wait()
	close(results)

	fmt.Printf("%d notes added\n", nCnt)
}

func createJNoteWorker(id int, wg *sync.WaitGroup, jnotes <-chan *jNote, results chan<- error) {
	defer wg.Done()

	for jn := range jnotes {
		if !utils.InSliceString(jn.Type, note.NoteTypes) {
			results <- errors.New("Invalid type")
			continue
		}

		book, err := dbConn.GetOrCreateBookByName(jn.Book)
		if err != nil {
			results <- err
			continue
		}

		tags := make(note.Tags, 0, len(jn.Tags))
		for _, t := range jn.Tags {
			tag, err := dbConn.GetOrCreateTagByName(t)
			if err != nil {
				results <- err
				continue
			}
			tags = append(tags, tag)
		}

		n := &note.Note{
			Created:  time.Now(),
			Modified: time.Now(),
			Book:     book,
			Type:     jn.Type,
			Title:    jn.Title,
			Body:     jn.Body,
			Tags:     tags,
		}

		err = saveNote(n)
		if err != nil {
			results <- err
			continue
		}

		fmt.Print("Note added: ")
		utils.PrintNoteColored(n, true)
		results <- nil
	}
}

func saveNote(n *note.Note) error {
	if err := dbConn.CreateNote(n); err != nil {
		return err
	}
	return idxConn.IndexNote(n)
}
