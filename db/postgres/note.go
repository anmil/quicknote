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
	"fmt"
	"strings"

	"github.com/anmil/quicknote"
)

// GetNoteByID returns the note for the given ID
func (d *Database) GetNoteByID(id int64) (*quicknote.Note, error) {
	sqlStr := `SELECT id, created, modified, bk_id, type, title, body FROM notes WHERE id = $1;`

	stmt, err := d.db.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	n := quicknote.NewNote()
	n.Book = quicknote.NewBook()

	err = stmt.QueryRow(id).Scan(&n.ID, &n.Created, &n.Modified, &n.Book.ID, &n.Type, &n.Title, &n.Body)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if err = d.LoadNoteTags(n); err != nil {
		return nil, err
	}

	if err = d.LoadBook(n.Book); err != nil {
		return nil, err
	}

	return n, nil
}

// GetNoteByNote Loads the note's ID, Created, and Modified fields
func (d *Database) GetNoteByNote(n *quicknote.Note) error {
	sqlStr := `SELECT id, created, modified FROM notes WHERE bk_id = $1 AND type = $2 AND title = $3 AND body = $4;`

	stmt, err := d.db.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(n.Book.ID, n.Type, n.Title, n.Body).
		Scan(&n.ID, &n.Created, &n.Modified)
	if err == sql.ErrNoRows {
		return nil
	}

	return err
}

// GetNotesByIDs returns all notes for the given Notebook
func (d *Database) GetNotesByIDs(ids []int64) (quicknote.Notes, error) {
	sqlStr := `SELECT id, created, modified, bk_id, type, title, body FROM notes WHERE id IN (%s);`

	// SQLite has a limit on the number of wild cards that can be given. We must split the query across multiple
	// calls if this number is exceeded. See splitSliceToChuck for more information
	chucks, err := splitSliceToChuck(ids, 1)
	if err != nil {
		return nil, err
	}

	notes := make(quicknote.Notes, 0, len(ids))
	for _, c := range chucks {
		cids := c.([]int64)

		// http://stackoverflow.com/questions/12990338/cannot-convert-string-to-interface
		qids := make([]interface{}, len(cids))
		pStr := make([]string, len(cids))
		for i, v := range cids {
			qids[i] = v
			pStr[i] = fmt.Sprintf("$%d", i+1)
		}
		query := fmt.Sprintf(sqlStr, strings.Join(pStr, ","))

		stmt, err := d.db.Prepare(query)
		if err != nil {
			return nil, err
		}
		defer stmt.Close()

		rows, err := stmt.Query(qids...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		ns, err := d.loadNotesFromRows(rows)
		if err != nil {
			return nil, err
		}
		notes = append(notes, ns...)
	}

	return notes, nil
}

// GetAllBookNotes returns all notes for the given Notebook
func (d *Database) GetAllBookNotes(book *quicknote.Book, sortBy, order string) (quicknote.Notes, error) {
	sqlStr := `SELECT id, created, modified, bk_id, type, title, body FROM notes WHERE bk_id = $1 ORDER BY %s %s;`

	// This would normally be a really bad idea (sql injection anyone?). But sortBy and order are taking
	// from command flags that are checked against a list of accepted values. The user is presented with
	// an error message if the give anything other than what is in those list.
	//
	// I have no other choice but to do this
	// See: http://stackoverflow.com/questions/30867337/golang-order-by-issue-with-mysql
	query := fmt.Sprintf(sqlStr, sortBy, order)

	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(book.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return d.loadNotesFromRows(rows)
}

// GetAllNotes returns all notes
func (d *Database) GetAllNotes(sortBy, order string) (quicknote.Notes, error) {
	sqlStr := `SELECT id, created, modified, bk_id, type, title, body FROM notes ORDER BY %s %s;`

	// See GetAllBookNotes for why I'm doing this
	query := fmt.Sprintf(sqlStr, sortBy, order)

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return d.loadNotesFromRows(rows)
}

// CreateNote saves the note to the database
func (d *Database) CreateNote(n *quicknote.Note) error {
	sqlStr := "INSERT INTO notes (created, modified, bk_id, type, title, body) " +
		"VALUES ($1,$2,$3,$4,$5,$6) RETURNING id;"

	tx, stmt, err := d.getTxStmt(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if err := stmt.QueryRow(n.Created, n.Modified, n.Book.ID, n.Type, n.Title, n.Body).Scan(&n.ID); err != nil {
		tx.Rollback()
		fmt.Println("Error 1")
		return err
	}

	if d.createTagRal(n, tx); err != nil {
		tx.Rollback()
		fmt.Println("Error 2")
		return err
	}

	return tx.Commit()
}

func (d *Database) EditNote(n *quicknote.Note) error {
	sqlStr := "UPDATE notes SET modified = $1, title = $2, body = $3 WHERE id = $4;"

	tx, stmt, err := d.getTxStmt(sqlStr)
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(n.Modified, n.Title, n.Body, n.ID); err != nil {
		tx.Rollback()
		return err
	}

	if d.deleteTagRal(n); err != nil {
		tx.Rollback()
		return err
	}
	if d.createTagRal(n, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// EditNoteByIDBook updates all notes for the given IDs with the Book bk's ID
func (d *Database) EditNoteByIDBook(ids []int64, bk *quicknote.Book) error {
	sqlStr1 := "UPDATE notes SET bk_id = $1 WHERE id in (%s);"
	sqlStr2 := "UPDATE note_book_tag SET bk_id = $1 WHERE note_id in (%s);"

	// SQLite has a limit on the number of wild cards that can be given. We must split the query across multiple
	// calls if this number is exceeded. See splitSliceToChuck for more information
	chucks, err := splitSliceToChuck(ids, 2)
	if err != nil {
		return err
	}

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	for _, c := range chucks {
		cids := c.([]int64)

		// http://stackoverflow.com/questions/12990338/cannot-convert-string-to-interface
		qids := make([]interface{}, len(cids))
		pStr := make([]string, len(cids))
		for i, v := range cids {
			qids[i] = v
			pStr[i] = fmt.Sprintf("$%d", i+2)
		}
		query1 := fmt.Sprintf(sqlStr1, strings.Join(pStr, ","))
		query2 := fmt.Sprintf(sqlStr2, strings.Join(pStr, ","))
		args := append([]interface{}{bk.ID}, qids...)

		stmt, err := tx.Prepare(query1)
		if err != nil {
			tx.Rollback()
			return err
		}
		if _, err = stmt.Exec(args...); err != nil {
			tx.Rollback()
			return err
		}

		stmt, err = tx.Prepare(query2)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = stmt.Exec(args...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// DeleteNote delete note from database
func (d *Database) DeleteNote(n *quicknote.Note) error {
	sqlStr := `DELETE FROM notes WHERE id = $1;`

	stmt, err := d.db.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(n.ID); err != nil {
		return err
	}

	return nil
}

func (d *Database) loadNotesFromRows(rows *sql.Rows) (quicknote.Notes, error) {
	books := make(map[int64]*quicknote.Book)
	notes := make(quicknote.Notes, 0)

	for rows.Next() {
		var bkID int64
		n := &quicknote.Note{}

		err := rows.Scan(&n.ID, &n.Created, &n.Modified, &bkID, &n.Type, &n.Title, &n.Body)
		if err != nil {
			return nil, err
		}

		if _, found := books[bkID]; !found {
			books[bkID] = &quicknote.Book{ID: bkID}
		}
		n.Book = books[bkID]

		if err = d.LoadNoteTags(n); err != nil {
			return nil, err
		}

		notes = append(notes, n)
	}

	for _, book := range books {
		if err := d.LoadBook(book); err != nil {
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}
