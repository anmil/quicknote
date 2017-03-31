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

package db

import (
	"errors"

	"github.com/anmil/quicknote/db/postgres"
	"github.com/anmil/quicknote/db/sqlite"
	"github.com/anmil/quicknote"
)

// ErrProviderNotSupported database provider given is not supported
var ErrProviderNotSupported = errors.New("Unsupported database provider")

// DB interface for the database providers
type DB interface {
	GetAllNotes(sortBy, order string) (quicknote.Notes, error)
	GetAllBookNotes(book *quicknote.Book, sortBy, order string) (quicknote.Notes, error)
	GetNoteByID(id int64) (*quicknote.Note, error)
	GetNoteByNote(n *quicknote.Note) error
	GetAllNotesByIDs(ids []int64) (quicknote.Notes, error)
	CreateNote(n *quicknote.Note) error
	EditNote(n *quicknote.Note) error
	DeleteNote(n *quicknote.Note) error

	GetAllBooks() (quicknote.Books, error)
	GetOrCreateBookByName(name string) (*quicknote.Book, error)
	GetBookByName(name string) (*quicknote.Book, error)
	CreateBook(b *quicknote.Book) error
	MergeBooks(b1 *quicknote.Book, b2 *quicknote.Book) error
	EditNoteByIDBook(ids []int64, bk *quicknote.Book) error
	EditBook(b1 *quicknote.Book) error
	LoadBook(b *quicknote.Book) error
	DeleteBook(bk *quicknote.Book) error

	GetAllBookTags(bk *quicknote.Book) (quicknote.Tags, error)
	GetAllTags() (quicknote.Tags, error)
	CreateTag(t *quicknote.Tag) error
	LoadNoteTags(n *quicknote.Note) error
	GetOrCreateTagByName(name string) (*quicknote.Tag, error)
	GetTagByName(name string) (*quicknote.Tag, error)

	Close() error
}

// NewDatabase returns a new database for the given provider
func NewDatabase(provider string, options ...string) (DB, error) {
	switch provider {
	case "sqlite":
		return sqlite.NewDatabase(options...)
	case "postgres":
		return postgres.NewDatabase(options...)
	default:
		return nil, ErrProviderNotSupported
	}
}
