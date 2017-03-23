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
	"fmt"
	"sort"
	"testing"

	"github.com/anmil/quicknote/test"
)

var bpText1 = `This is test 1 of the basic parser
#basic #test #parser

Lorem ipsum dolor sit amet, consectetur adipiscing elit.
Nulla tincidunt diam eu purus laoreet condimentum. Duis
tempus, turpis vitae varius#ullamcorper, sapien erat
cursus lacus, et lacinia ligula dolor quis nibh.`

var bpText1Title = "This is test 1 of the basic parser"
var bpText1Tags = []string{"basic", "parser", "test"}
var bpText1Body = `#basic #test #parser

Lorem ipsum dolor sit amet, consectetur adipiscing elit.
Nulla tincidunt diam eu purus laoreet condimentum. Duis
tempus, turpis vitae varius#ullamcorper, sapien erat
cursus lacus, et lacinia ligula dolor quis nibh.`

var bpText2 = `This is #test 2 of the #basic #parser

Lorem ipsum dolor sit amet, consectetur adipiscing elit.
Nulla tincidunt diam eu purus laoreet condimentum. Duis
tempus, turpis vitae varius ullamcorper, sapien erat
cursus lacus, et lacinia ligula dolor #quis nibh.#`

var bpText2Title = "This is #test 2 of the #basic #parser"
var bpText2Tags = []string{"basic", "parser", "quis", "test"}
var bpText2Body = `Lorem ipsum dolor sit amet, consectetur adipiscing elit.
Nulla tincidunt diam eu purus laoreet condimentum. Duis
tempus, turpis vitae varius ullamcorper, sapien erat
cursus lacus, et lacinia ligula dolor #quis nibh.#`

func TestBasicParserUnit(t *testing.T) {
	parser := &BasicParser{}
	parser.Parse(bpText1)

	if parser.Title() != bpText1Title {
		t.Error("Parser returned incorrect title")
	}
	tags := parser.Tags()
	sort.Strings(tags)
	if !test.StringSliceEq(tags, bpText1Tags) {
		t.Error("Parser returned incorrect tags")
	}
	if parser.Body() != bpText1Body {
		t.Error("Parser returned incorrect body")
	}

	parser = &BasicParser{}
	parser.Parse(bpText2)

	if parser.Title() != bpText2Title {
		t.Error("Parser returned incorrect title")
	}
	tags = parser.Tags()
	sort.Strings(tags)
	if !test.StringSliceEq(tags, bpText2Tags) {
		fmt.Println(parser.Tags())
		fmt.Println(bpText2Tags)
		t.Error("Parser returned incorrect tags")
	}
	if parser.Body() != bpText2Body {
		t.Error("Parser returned incorrect body")
	}
}
