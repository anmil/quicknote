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

package parser

import (
	"strings"
	"unicode"
)

// BasicParser Is the default parser for notes. The first sentence
// is used as the title. Everything after the first sentence is
// used as the note's body. Any word starting with `#` is parsed as
// a tag.
type BasicParser struct {
	title string
	tags  []string
	body  string
}

// Title returns the parsed title
func (p *BasicParser) Title() string {
	return p.title
}

// Tags returns the parsed tags
func (p *BasicParser) Tags() []string {
	return p.tags
}

// Body returns the parsed body
func (p *BasicParser) Body() string {
	return p.body
}

// Parse parses the text for the note's title, tags, and body
func (p *BasicParser) Parse(text string) {
	p.title, p.body = splitTitleBody(text)
	p.tags = getTags(text)
}

func splitTitleBody(text string) (string, string) {
	// If there is only one line, than we just have the title
	title := text
	body := ""

	// First line is treated as the title and the rest
	// is the body
	parts := strings.SplitN(text, "\n", 2)
	title = parts[0]
	if len(parts) == 2 {
		body = parts[1]
	}

	return strings.TrimSpace(title), strings.TrimSpace(body)
}

func getTags(text string) []string {
	tagCount := 0

	// Using a map as a makeshift Set
	tags := make(map[string]bool)

	tStart := 0
	tEnd := 0

	// A tag is any word that starts with the '#'
	// character. Any '#' character that does not have
	// a whitespace preceding it is ignored except when
	// its the first character of the string.
	//
	// Duplicate tags are ignored.
	for i := 0; i < len(text); {
		tStart = nextTagIndex(text, tEnd)
		if tStart == -1 {
			break
		}

		tEnd = getTagEndIndex(text, tStart+1)
		tag := text[tStart+1 : tEnd]

		if unicode.IsPunct(rune(tag[len(tag)-1])) {
			tag = tag[:len(tag)-1]
		}

		if _, found := tags[tag]; !found && len(tag) > 1 {
			tags[strings.ToLower(tag)] = true
			tagCount++
		}

		i = tEnd
	}

	// Get all the map's keys, which are our tags
	keys := make([]string, 0, tagCount)
	for k := range tags {
		keys = append(keys, k)
	}
	return keys
}

func nextTagIndex(text string, start int) int {
	for i := start; i < len(text); i++ {
		if i == 0 && isTag(text[i]) {
			return i
		}
		if i > 0 && (unicode.IsSpace(rune(text[i-1])) && isTag(text[i])) {
			return i
		}
	}
	return -1
}

func getTagEndIndex(text string, start int) int {
	i := start
	for ; i < len(text); i++ {
		if text[i] == ',' || unicode.IsSpace(rune(text[i])) {
			return i
		}
	}
	return i
}

func isTag(c byte) bool {
	return c == '#'
}
