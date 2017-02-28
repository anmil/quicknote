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

package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/anmil/quicknote/note"
)

// GetAllBooks returns all Books
func (d *Database) GetAllBooks() ([]*note.Book, error) {
	sqlStr := "SELECT id, created, modified, name FROM books;"

	rows, err := d.db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := make([]*note.Book, 0)
	for rows.Next() {
		b := note.NewBook()
		err := rows.Scan(&b.ID, &b.Created, &b.Modified, &b.Name)
		if err != nil {
			return nil, err
		}

		books = append(books, b)
	}

	return books, nil
}

// GetOrCreateBookByName gets the Book by name creating it if it does not exists
func (d *Database) GetOrCreateBookByName(name string) (*note.Book, error) {
	if len(name) == 0 {
		return nil, errors.New("No Notebook name given")
	}

	bk, err := d.GetBookByName(name)
	if err != nil {
		return nil, err
	}
	if bk == nil {
		bk = &note.Book{
			Created:  time.Now(),
			Modified: time.Now(),
			Name:     name,
		}
		err = d.CreateBook(bk)
		if err != nil {
			return nil, err
		}
	}

	return bk, nil
}

// GetBookByName returns the Book for the given name
func (d *Database) GetBookByName(name string) (*note.Book, error) {
	sqlStr := "SELECT id, created, modified, name FROM books WHERE name = $1;"

	stmt, err := d.db.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	b := note.NewBook()
	err = stmt.QueryRow(name).Scan(&b.ID, &b.Created, &b.Modified, &b.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return b, nil
}

// LoadBook loads the Note's Book
func (d *Database) LoadBook(b *note.Book) error {
	sqlStr := "SELECT created, modified, name FROM books WHERE id = $1;"

	stmt, err := d.db.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if err = stmt.QueryRow(b.ID).Scan(&b.Created, &b.Modified, &b.Name); err != nil {
		return err
	}

	return nil
}

// CreateBook saves the Book to the database
func (d *Database) CreateBook(b *note.Book) error {
	sqlStr := "INSERT INTO books (created, modified, name) VALUES ($1,$2,$3) RETURNING id;"

	tx, stmt, err := d.getTxStmt(sqlStr)
	if err != nil {
		return err
	}

	if err := stmt.QueryRow(b.Created, b.Modified, b.Name).Scan(&b.ID); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// MergeBooks merge all notes from Book b1 into Book b2
func (d *Database) MergeBooks(b1 *note.Book, b2 *note.Book) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	sqlStr := "UPDATE notes SET bk_id = $1, modified = $2 WHERE bk_id = $3;"
	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec(b2.ID, time.Now(), b1.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	sqlStr = "UPDATE note_book_tag SET bk_id = $1 WHERE bk_id = $2;"
	stmt, err = tx.Prepare(sqlStr)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec(b2.ID, b1.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	sqlStr = "DELETE FROM books WHERE id = $1;"
	stmt, err = tx.Prepare(sqlStr)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec(b1.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// EditBook change the book name
func (d *Database) EditBook(b *note.Book) error {
	sqlStr := "UPDATE books SET name = $1, modified = $2 where id = $3;"

	tx, stmt, err := d.getTxStmt(sqlStr)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(b.Name, time.Now(), b.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// DeleteBook deletes the Book from the database
func (d *Database) DeleteBook(bk *note.Book) error {
	sqlStr := "DELETE FROM books WHERE id = $1;"

	stmt, err := d.db.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(bk.ID); err != nil {
		return err
	}

	return nil
}
