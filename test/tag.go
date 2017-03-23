package test

import (
	"sort"
	"testing"

	"github.com/anmil/quicknote/note"
)

var AllTags note.Tags
var noteTags map[string]*note.Tag

func getTag(name string) *note.Tag {
	if t, found := noteTags[name]; found {
		return t
	}

	t := note.NewTag()
	t.Name = name
	noteTags[name] = t
	AllTags = append(AllTags, t)

	return t
}

func CheckTags(t *testing.T, tag1, tag2 note.Tags) {
	nnTags := note.Tags{}
	for _, t := range tag1 {
		nnTags = append(nnTags, t)
	}
	sort.Sort(nnTags)

	nTags := note.Tags{}
	for _, t := range tag2 {
		nTags = append(nTags, t)
	}
	sort.Sort(nTags)

	if !TagSliceEq(nnTags, nTags) {
		t.Fatal("Did not received the corrected Tags")
	}
}

func TagSliceEq(a, b note.Tags) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Name != b[i].Name {
			return false
		}
	}
	return true
}
