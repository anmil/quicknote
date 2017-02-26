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

	"github.com/anmil/quicknote/db/sqlite"
	"github.com/anmil/quicknote/note"
)

// ErrProviderNotSupported database provider given is not supported
var ErrProviderNotSupported = errors.New("Unsupported database provider")

// DB interface for the database providers
type DB interface {
	GetAllNotes(sortBy, order string) ([]*note.Note, error)
	GetAllBookNotes(book *note.Book, sortBy, order string) ([]*note.Note, error)
	GetNoteByID(id int64) (*note.Note, error)
	GetAllNotesByIDs(ids []int64) ([]*note.Note, error)
	CreateNote(n *note.Note) error
	EditNote(n *note.Note) error
	DeleteNote(n *note.Note) error

	GetAllBooks() ([]*note.Book, error)
	GetOrCreateBookByName(name string) (*note.Book, error)
	GetBookByName(name string) (*note.Book, error)
	CreateBook(b *note.Book) error
	MergeBooks(b1 *note.Book, b2 *note.Book) error
	EditNoteByIDBook(ids []int64, bk *note.Book) error
	EditBook(b1 *note.Book) error
	LoadBook(b *note.Book) error
	DeleteBook(bk *note.Book) error

	GetAllBookTags(bk *note.Book) ([]*note.Tag, error)
	GetAllTags() ([]*note.Tag, error)
	CreateTag(t *note.Tag) error
	LoadNoteTags(n *note.Note) error
	GetOrCreateTagByName(name string) (*note.Tag, error)
	GetTagByName(name string) (*note.Tag, error)

	Close()
}

// NewDatabase returns a new database for the given provider
func NewDatabase(provider string, options ...string) (DB, error) {
	switch provider {
	case "sqlite":
		return sqlite.NewDatabase(options...)
	default:
		return nil, ErrProviderNotSupported
	}
}
