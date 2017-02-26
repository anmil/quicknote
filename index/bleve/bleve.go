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

package bleve

import (
	"fmt"
	"strconv"
	"time"

	"github.com/blevesearch/bleve"

	"github.com/anmil/quicknote/note"
)

type indexNote struct {
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Type     string    `json:"type"`
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	Book     string    `json:"book"`
	Tags     []string  `json:"tags"`
}

// Index provides the interface to Bleve
type Index struct {
	db bleve.Index
}

// NewIndex returns a new Index
func NewIndex(indexPath string) (*Index, error) {
	indexMapping := bleve.NewIndexMapping()
	index, err := bleve.New(indexPath, indexMapping)
	if err == bleve.ErrorIndexPathExists {
		index, err = bleve.Open(indexPath)
	}
	if err != nil {
		return nil, err
	}
	return &Index{db: index}, nil
}

// IndexNote creates or updates a note in Bleve index
func (b *Index) IndexNote(n *note.Note) error {
	iN := &indexNote{
		Created:  n.Created,
		Modified: n.Modified,
		Type:     n.Type,
		Title:    n.Title,
		Body:     n.Body,
		Book:     n.Book.Name,
		Tags:     n.GetTagStringArray(),
	}

	err := b.db.Index(strconv.FormatInt(n.ID, 10), iN)
	return err
}

// IndexNotes creates or updates a list of notes in Bleve index
func (b *Index) IndexNotes(notes []*note.Note) error {
	for _, n := range notes {
		if err := b.IndexNote(n); err != nil {
			return err
		}
	}
	return nil
}

// SearchNote sends a search query to Bleve using QueryStringQuery
func (b *Index) SearchNote(query string, limit, offset int) ([]int64, uint64, error) {
	q := bleve.NewQueryStringQuery(query)
	search := bleve.NewSearchRequest(q)
	search.Size = limit
	search.From = offset
	res, err := b.db.Search(search)
	if err != nil {
		return nil, 0, err
	}

	ids := make([]int64, 0, len(res.Hits))
	for _, h := range res.Hits {
		id, err := strconv.ParseInt(h.ID, 10, 64)
		if err != nil {
			return nil, 0, err
		}
		ids = append(ids, id)
	}

	return ids, res.Total, err
}

// SearchNotePhrase sends a search query to Bleve using Prefix query
// If bk is given, only notes for that Book are queried.
func (b *Index) SearchNotePhrase(query string, bk *note.Book, sort string, limit, offset int) ([]int64, uint64, error) {
	boolQuery := bleve.NewBooleanQuery()

	matchPrefixQuery := bleve.NewPrefixQuery(query)
	boolQuery.AddMust(matchPrefixQuery)

	if bk != nil {
		matchBookQuery := bleve.NewQueryStringQuery(fmt.Sprintf("+book:%s", bk.Name))
		boolQuery.AddMust(matchBookQuery)
	}

	search := bleve.NewSearchRequest(boolQuery)
	search.Size = limit
	search.From = offset

	res, err := b.db.Search(search)
	if err != nil {
		return nil, 0, err
	}

	ids := make([]int64, 0, len(res.Hits))
	for _, h := range res.Hits {
		id, err := strconv.ParseInt(h.ID, 10, 64)
		if err != nil {
			return nil, 0, err
		}
		ids = append(ids, id)
	}

	// Go has no build in function for reversing a array, and Bleve does not have
	// a simple way either. You can sort by field, but I could not figure out how
	// to just reverse results
	if sort == "asc" {
		for i := 0; i < len(ids)/2; i++ {
			j := len(ids) - i - 1
			ids[i], ids[j] = ids[j], ids[i]
		}
	}

	return ids, res.Total, err
}

// DeleteNote deletes note from index
func (b *Index) DeleteNote(n *note.Note) error {
	return b.db.Delete(strconv.FormatInt(n.ID, 10))
}

// DeleteBook deletes all notes in the index for the notebook
func (b *Index) DeleteBook(bk *note.Book) error {
	batch, ids, err := b.getNextDeleteBatch(bk)
	if err != nil {
		return err
	}

	for len(ids) > 0 {
		for _, id := range ids {
			batch.Delete(id)
		}
		b.db.Batch(batch)
		batch, ids, err = b.getNextDeleteBatch(bk)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Index) getNextDeleteBatch(bk *note.Book) (*bleve.Batch, []string, error) {
	batch := b.db.NewBatch()

	query := fmt.Sprintf("+Book:%s", bk.Name)
	q := bleve.NewQueryStringQuery(query)
	search := bleve.NewSearchRequest(q)
	search.Size = 1000
	res, err := b.db.Search(search)
	if err != nil {
		return nil, nil, err
	}

	ids := make([]string, 0, res.Total)
	for _, h := range res.Hits {
		ids = append(ids, h.ID)
	}

	return batch, ids, err
}
