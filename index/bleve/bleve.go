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
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/blevesearch/bleve"
	bquery "github.com/blevesearch/bleve/search/query"

	"github.com/anmil/quicknote/note"
)

type indexNote struct {
	ID       int64     `json:"id"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Type     string    `json:"type"`
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	Book     string    `json:"book"`
	Tags     []string  `json:"tags"`
}

type bIndex struct {
	Index bleve.Index
	mux   sync.Mutex
}

func (b *bIndex) BatchIndex(wg *sync.WaitGroup, notes []*indexNote) {
	b.mux.Lock()
	defer b.mux.Unlock()
	defer wg.Done()

	batch := b.Index.NewBatch()
	for _, iN := range notes {
		err := batch.Index(strconv.FormatInt(iN.ID, 10), iN)
		if err != nil {
			panic(err)
		}
	}

	err := b.Index.Batch(batch)
	if err != nil {
		panic(err)
	}
}

func (b *bIndex) DeleteBook(wg *sync.WaitGroup, bk *note.Book) {
	defer wg.Done()

	ids, err := b.getNextDeleteBatch(bk)
	if err != nil {
		panic(err)
	}

	for len(ids) > 0 {
		batch := b.Index.NewBatch()
		for _, id := range ids {
			batch.Delete(id)
		}

		err := b.Index.Batch(batch)
		if err != nil {
			panic(err)
		}

		ids, err = b.getNextDeleteBatch(bk)
		if err != nil {
			panic(err)
		}
	}
}

func (b *bIndex) getNextDeleteBatch(bk *note.Book) ([]string, error) {
	query := fmt.Sprintf("+book:%s", bk.Name)
	q := bleve.NewQueryStringQuery(query)

	search := bleve.NewSearchRequest(q)
	search.Size = 1000

	res, err := b.Index.Search(search)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, res.Total)
	for _, h := range res.Hits {
		ids = append(ids, h.ID)
	}

	return ids, err
}

// Index provides the interface to Bleve
type Index struct {
	db bleve.IndexAlias

	shards  int
	indexes []*bIndex

	indexIdx     int
	indexIdxFile string
}

// NewIndex returns a new Index
func NewIndex(indexPath string, shards int) (*Index, error) {
	indexMapping := bleve.NewIndexMapping()

	bindexes := make([]*bIndex, shards)
	indexes := make([]bleve.Index, shards)

	for i := 0; i < shards; i++ {
		p := path.Join(indexPath, fmt.Sprintf("index-%d.bleve", i))
		index, err := bleve.New(p, indexMapping)
		if err == bleve.ErrorIndexPathExists {
			index, err = bleve.Open(p)
			if err != nil {
				return nil, err
			}
		}
		indexes[i] = index
		bindexes[i] = &bIndex{Index: index}
	}

	indexAlias := bleve.NewIndexAlias(indexes...)
	indexIdxFile := path.Join(indexPath, "current_index")

	idx := &Index{
		db:           indexAlias,
		shards:       shards,
		indexes:      bindexes,
		indexIdxFile: indexIdxFile,
	}

	err := idx.loadIndexIdx()
	if err != nil {
		return nil, err
	}

	return idx, nil
}

func (b *Index) getIndex() *bIndex {
	b.indexIdx++
	if b.indexIdx >= b.shards {
		b.indexIdx = 0
	}
	err := b.saveIndexIdx()
	if err != nil {
		panic(err)
	}
	return b.indexes[b.indexIdx]
}

func (b *Index) getDocIndex(id string) (*bIndex, error) {
	for _, idx := range b.indexes {
		if doc, err := idx.Index.Document(id); err != nil {
			return nil, err
		} else if doc != nil {
			return idx, nil
		}
	}
	return nil, nil
}

func (b *Index) loadIndexIdx() error {
	if _, err := os.Stat(b.indexIdxFile); !os.IsNotExist(err) {
		data, err := ioutil.ReadFile(b.indexIdxFile)
		if err != nil {
			return nil
		}
		b.indexIdx, err = strconv.Atoi(string(data))
		return err
	}

	b.indexIdx = 0
	return nil
}

func (b *Index) saveIndexIdx() error {
	s := strconv.Itoa(b.indexIdx)
	return ioutil.WriteFile(b.indexIdxFile, []byte(s), 0600)
}

// IndexNote creates or updates a note in Bleve index
func (b *Index) IndexNote(n *note.Note) error {
	iN := &indexNote{
		ID:       n.ID,
		Created:  n.Created,
		Modified: n.Modified,
		Type:     n.Type,
		Title:    n.Title,
		Body:     n.Body,
		Book:     n.Book.Name,
		Tags:     n.GetTagStringArray(),
	}

	idS := strconv.FormatInt(n.ID, 10)
	idx, err := b.getDocIndex(idS)
	if err != nil {
		return err
	}

	if idx != nil {
		err = idx.Index.Index(idS, iN)
	} else {
		err = b.getIndex().Index.Index(idS, iN)
	}

	return err
}

// IndexNotes creates or updates a list of notes in Bleve index
func (b *Index) IndexNotes(notes note.Notes) error {
	var wg sync.WaitGroup
	var batchSize int64

	index := b.getIndex()
	bNotes := make([]*indexNote, 0)
	for _, n := range notes {
		iN := &indexNote{
			ID:       n.ID,
			Created:  n.Created,
			Modified: n.Modified,
			Type:     n.Type,
			Title:    n.Title,
			Body:     n.Body,
			Book:     n.Book.Name,
			Tags:     n.GetTagStringArray(),
		}

		// Check if this note has been indexed already
		// update it if so
		idS := strconv.FormatInt(n.ID, 10)
		if idx, err := b.getDocIndex(idS); err != nil {
			return err
		} else if idx != nil {
			if err = idx.Index.Index(idS, iN); err != nil {
				return err
			}
			continue
		}

		batchSize += 16 + // timestamps are 8 bytes each
			int64(len(n.Type)) +
			int64(len(n.Title)) +
			int64(len(n.Body)) +
			int64(len(n.Book.Name)) +
			int64(len(n.GetTagStringArray()))

		bNotes = append(bNotes, iN)

		// After some experimenting I've found that batch performance depends more
		// on the byte size than the number of records. Using Bolt as the store
		// provider. Some testing shows that around 64KB is the sweet spot.
		if batchSize >= 65536 {
			wg.Add(1)
			go index.BatchIndex(&wg, bNotes)
			bNotes = make([]*indexNote, 0)

			index = b.getIndex()
			batchSize = 0
		}
	}

	if len(bNotes) > 0 {
		wg.Add(1)
		go index.BatchIndex(&wg, bNotes)
	}

	wg.Wait()

	nCnt := uint64(0)
	for _, idx := range b.indexes {
		cnt, _ := idx.Index.DocCount()
		nCnt += cnt
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

	// Bleve does not support phrase prefix query natively
	// https://github.com/blevesearch/bleve/issues/377
	var disquery bquery.Query
	words := strings.Fields(query)
	if len(words) == 1 {
		disquery = bleve.NewDisjunctionQuery(
			bleve.NewPrefixQuery(query),
			bleve.NewMatchQuery(query),
		)
	} else {
		var phrase string
		if len(words) == 2 {
			phrase = words[0]
		} else {
			phrase = strings.Join(words[0:len(words)-2], " ")
		}

		disquery = bleve.NewDisjunctionQuery(
			bleve.NewMatchPhraseQuery(query), // whole thing as a phrase, or..
			bleve.NewConjunctionQuery( // phrase + prefix
				bleve.NewMatchPhraseQuery(phrase),
				bleve.NewPrefixQuery(words[len(words)-1]),
			),
		)
	}
	boolQuery.AddMust(disquery)

	// matchPrefixQuery := bleve.NewPrefixQuery(query)
	// boolQuery.AddMust(matchPrefixQuery)

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
	for _, i := range b.indexes {
		err := i.Index.Delete(strconv.FormatInt(n.ID, 10))
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteBook deletes all notes in the index for the notebook
func (b *Index) DeleteBook(bk *note.Book) error {
	var wg sync.WaitGroup

	for _, i := range b.indexes {
		wg.Add(1)
		go i.DeleteBook(&wg, bk)
	}

	wg.Wait()
	return nil
}
