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

package elastic

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/anmil/quicknote/note"

	"golang.org/x/net/context"
	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	// TitleBoost boost value for the title field
	TitleBoost = 0.8

	// TagsBoost boost value for the tags field
	TagsBoost = 0.6

	// BodyBoost boost value for the body field
	BodyBoost = 0.5

	// MaxExpansions max expansion sets the limit on how many
	// documents the prefix max will match before returning. Prefix
	// matching is resource intensive.
	//
	// For more details see ElasticSearch's docs on "Query-Time Search-as-You-Type"
	MaxExpansions = 50

	// Slop how much slop to give when matching the order and position of the words
	// For more details see ElasticSearch's docs on "Query-Time Search-as-You-Type"
	Slop = 20
)

// Index provides the interface to ElasticSearch
type Index struct {
	client    *elastic.Client
	indexName string
}

// NewIndex returns a new Index
func NewIndex(host, idxName string) (*Index, error) {
	ctx := context.Background()

	var err error
	var client *elastic.Client
	if len(host) != 0 {
		client, err = elastic.NewClient(elastic.SetURL(host))
	} else {
		client, err = elastic.NewClient()
	}

	if err != nil {
		return nil, err
	}

	// Make sure our index exists
	exists, err := client.IndexExists(idxName).Do(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		createIndex, err := client.CreateIndex(idxName).Do(ctx)
		if err != nil {
			return nil, err
		}
		if !createIndex.Acknowledged {
			return nil, errors.New("ElasticSearch failed to acknowledged the new index")
		}
	}

	return &Index{client: client, indexName: idxName}, nil
}

// IndexNote creates or updates a note in ElasticSearch index
func (b *Index) IndexNote(n *note.Note) error {
	ctx := context.Background()

	_, err := b.client.Index().
		Index(b.indexName).
		Type("note").
		Id(strconv.FormatInt(n.ID, 10)).
		BodyJson(n).
		Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

// IndexNotes creates or updates a list of notes in ElasticSearch index
func (b *Index) IndexNotes(notes []*note.Note) error {
	for _, n := range notes {
		b.IndexNote(n)
	}
	return nil
}

// SearchNote sends a search query to ElasticSearch using QueryStringQuery
func (b *Index) SearchNote(query string, limit, offset int) ([]int64, uint64, error) {
	ctx := context.Background()

	stringQuery := elastic.NewQueryStringQuery(query)
	stringQuery.FieldWithBoost("title", TitleBoost)
	stringQuery.FieldWithBoost("tags", TagsBoost)
	stringQuery.FieldWithBoost("body", BodyBoost)

	scoreSort := elastic.NewScoreSort()
	scoreSort.Asc()

	searchResult, err := b.client.Search().
		Index(b.indexName).
		SortBy(scoreSort).
		Query(stringQuery).
		From(offset).Size(limit).
		Do(ctx)
	if err != nil {
		return nil, 0, err
	}

	return b.getNoteIDsFromResults(searchResult)
}

// SearchNotePhrase sends a search query to ElasticSearch using Phrase Prefix query
// If bk is given, only notes for that Book are queried.
func (b *Index) SearchNotePhrase(query string, bk *note.Book, sort string, limit, offset int) ([]int64, uint64, error) {
	ctx := context.Background()

	matchPhrasePrefixQuery := elastic.NewMultiMatchQuery(query)
	matchPhrasePrefixQuery.Type("phrase_prefix")
	matchPhrasePrefixQuery.FieldWithBoost("title", TitleBoost)
	matchPhrasePrefixQuery.FieldWithBoost("tags", TagsBoost)
	matchPhrasePrefixQuery.FieldWithBoost("body", BodyBoost)
	matchPhrasePrefixQuery.MaxExpansions(MaxExpansions)
	matchPhrasePrefixQuery.Slop(Slop)

	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(matchPhrasePrefixQuery)

	if bk != nil {
		notebookMatchQuery := elastic.NewMatchQuery("book", bk.Name)
		boolQuery.Must(notebookMatchQuery)
	}

	scoreSort := elastic.NewScoreSort()

	if sort == "asc" {
		scoreSort.Asc()
	}

	searchResult, err := b.client.Search().
		Index(b.indexName).
		SortBy(scoreSort).
		Query(boolQuery).
		From(offset).Size(limit).
		Do(ctx)
	if err != nil {
		return nil, 0, err
	}

	return b.getNoteIDsFromResults(searchResult)
}

func (b *Index) getNoteIDsFromResults(sr *elastic.SearchResult) ([]int64, uint64, error) {
	total := uint64(sr.Hits.TotalHits)
	ids := make([]int64, 0, len(sr.Hits.Hits))
	for _, h := range sr.Hits.Hits {
		id, err := strconv.ParseInt(h.Id, 10, 64)
		if err != nil {
			return nil, 0, err
		}
		ids = append(ids, id)
	}

	return ids, total, nil
}

// DeleteNote deletes note from index
func (b *Index) DeleteNote(n *note.Note) error {
	ctx := context.Background()

	_, err := b.client.Delete().
		Index(b.indexName).
		Type("note").
		Id(strconv.FormatInt(n.ID, 10)).
		Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

// DeleteBook deletes all notes in the index for the notebook
func (b *Index) DeleteBook(bk *note.Book) error {
	ctx := context.Background()

	query := fmt.Sprintf("book:%s", bk.Name)
	deleteQuery := elastic.NewQueryStringQuery(query)
	_, err := b.client.DeleteByQuery(b.indexName).
		Type("note").
		Query(deleteQuery).
		Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

// DeleteIndex deletes this index
func (b *Index) DeleteIndex() error {
	ctx := context.Background()

	deleteIndex, err := b.client.DeleteIndex(b.indexName).Do(ctx)
	if err != nil {
		return err
	}
	if !deleteIndex.Acknowledged {
		return errors.New("Delete Index was not acknowledged")
	}
	return nil
}
