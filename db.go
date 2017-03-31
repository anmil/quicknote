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

package quicknote

// DB interface for the database providers
type DB interface {
	GetAllNotes(sortBy, order string) (Notes, error)
	GetAllBookNotes(book *Book, sortBy, order string) (Notes, error)
	GetNoteByID(id int64) (*Note, error)
	GetNoteByNote(n *Note) error
	GetNotesByIDs(ids []int64) (Notes, error)
	CreateNote(n *Note) error
	EditNote(n *Note) error
	DeleteNote(n *Note) error

	GetAllBooks() (Books, error)
	GetOrCreateBookByName(name string) (*Book, error)
	GetBookByName(name string) (*Book, error)
	CreateBook(b *Book) error
	MergeBooks(b1 *Book, b2 *Book) error
	EditNoteByIDBook(ids []int64, bk *Book) error
	EditBook(b1 *Book) error
	LoadBook(b *Book) error
	DeleteBook(bk *Book) error

	GetAllBookTags(bk *Book) (Tags, error)
	GetAllTags() (Tags, error)
	CreateTag(t *Tag) error
	LoadNoteTags(n *Note) error
	GetOrCreateTagByName(name string) (*Tag, error)
	GetTagByName(name string) (*Tag, error)

	Close() error
}
