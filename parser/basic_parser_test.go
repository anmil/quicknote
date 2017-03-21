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
