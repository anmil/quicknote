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
	"strings"
	"time"

	"github.com/anmil/quicknote/note"
)

// GetAllBookTags returns all tags for the given Book
func (d *Database) GetAllBookTags(bk *note.Book) ([]*note.Tag, error) {
	sqlStr := "SELECT id, created, modified, name FROM tags WHERE id in " +
		"(SELECT tag_id FROM note_book_tag WHERE bk_id = $1);"

	stmt, err := d.db.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(bk.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return d.loadTagsFromRows(rows)
}

// GetAllTags returns all tags
func (d *Database) GetAllTags() ([]*note.Tag, error) {
	sqlStr := "SELECT id, created, modified, name FROM tags;"

	rows, err := d.db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return d.loadTagsFromRows(rows)
}

// GetOrCreateTagByName returns a tag, creating it if it does not exists
func (d *Database) GetOrCreateTagByName(name string) (*note.Tag, error) {
	if len(name) == 0 {
		return nil, errors.New("No Tag name given")
	}

	tg, err := d.GetTagByName(name)
	if err != nil {
		return nil, err
	}
	if tg == nil {
		tg = &note.Tag{
			Created:  time.Now(),
			Modified: time.Now(),
			Name:     name,
		}
		err = d.CreateTag(tg)
		if err != nil {
			return nil, err
		}
	}

	return tg, nil
}

// GetTagByName returns the tag with the given name
func (d *Database) GetTagByName(name string) (*note.Tag, error) {
	sqlStr := "SELECT id, created, modified, name FROM tags WHERE name = $1;"

	stmt, err := d.db.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	t := note.NewTag()
	err = stmt.QueryRow(name).Scan(&t.ID, &t.Created, &t.Modified, &t.Name)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return t, nil
}

// LoadNoteTags loads all the tags for the given Note
func (d *Database) LoadNoteTags(n *note.Note) error {
	sqlStr := "SELECT id, created, modified, name FROM tags WHERE id in " +
		"(SELECT tag_id FROM note_tag WHERE note_id = $1);"

	stmt, err := d.db.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(n.ID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	n.Tags, err = d.loadTagsFromRows(rows)
	return err
}

// CreateTag saves the tag to the database
func (d *Database) CreateTag(t *note.Tag) error {
	sqlStr := "INSERT INTO tags (created, modified, name) VALUES ($1,$2,$3) RETURNING id;"

	tx, stmt, err := d.getTxStmt(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if err := stmt.QueryRow(t.Created, t.Modified, t.Name).Scan(&t.ID); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// GetTagsByName returns the Tag for the given name
func (d *Database) GetTagsByName(name string) (*note.Tag, error) {
	return nil, nil
}

func (d *Database) loadTagsFromRows(rows *sql.Rows) ([]*note.Tag, error) {
	tags := make([]*note.Tag, 0)
	for rows.Next() {
		t := note.NewTag()
		err := rows.Scan(&t.ID, &t.Created, &t.Modified, &t.Name)
		if err != nil {
			return nil, err
		}

		tags = append(tags, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (d *Database) createTagRal(n *note.Note, tx *sql.Tx) error {
	for _, t := range n.Tags {
		if err := d.createNoteTagRel(n, t, tx); err != nil {
			return err
		}
		if err := d.createNoteBookTagRel(n, t, tx); err != nil {
			return err
		}
	}
	return nil
}

func (d *Database) createNoteTagRel(n *note.Note, t *note.Tag, tx *sql.Tx) error {
	sqlStr := "INSERT INTO note_tag (note_id, tag_id) VALUES ($1,$2);"

	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(n.ID, t.ID)
	if err != nil && !strings.Contains(err.Error(), "violates unique constraint") {
		return err
	}

	return nil
}

func (d *Database) createNoteBookTagRel(n *note.Note, t *note.Tag, tx *sql.Tx) error {
	sqlStr := "INSERT INTO note_book_tag (note_id, bk_id, tag_id) VALUES ($1,$2,$3);"

	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(n.ID, n.Book.ID, t.ID)
	if err != nil && !strings.Contains(err.Error(), "violates unique constraint") {
		return err
	}

	return nil
}

func (d *Database) deleteTagRal(n *note.Note) error {
	if err := d.deleteNoteTagsRel(n); err != nil {
		return err
	}
	if err := d.deleteNoteNookTagsRel(n); err != nil {
		return err
	}
	return nil
}

func (d *Database) deleteNoteTagsRel(n *note.Note) error {
	sqlStr := "DELETE FROM note_tag WHERE note_id = $1"

	stmt, err := d.db.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(n.ID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	return nil
}

func (d *Database) deleteNoteNookTagsRel(n *note.Note) error {
	sqlStr := "DELETE FROM note_book_tag WHERE note_id = $1"

	stmt, err := d.db.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(n.ID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	return nil
}
