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

package sqlite

import (
	"database/sql"
	"strings"
	"time"

	"github.com/anmil/quicknote"
)

// GetAllBookTags returns all tags for the given Book
func (d *Database) GetAllBookTags(bk *quicknote.Book) (quicknote.Tags, error) {
	d.mux.Lock()
	defer d.mux.Unlock()

	sqlStr := "SELECT id, created, modified, name FROM tags WHERE id in " +
		"(SELECT tag_id FROM note_book_tag WHERE bk_id = ?);"

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
func (d *Database) GetAllTags() (quicknote.Tags, error) {
	d.mux.Lock()
	defer d.mux.Unlock()

	sqlStr := "SELECT id, created, modified, name FROM tags;"

	rows, err := d.db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return d.loadTagsFromRows(rows)
}

// GetOrCreateTagByName returns a tag, creating it if it does not exists
func (d *Database) GetOrCreateTagByName(name string) (*quicknote.Tag, error) {
	if t := d.getFromTagCache(name); t != nil {
		return t, nil
	}

	t, err := d.GetTagByName(name)
	if err != nil {
		return nil, err
	}
	if t == nil {
		t = &quicknote.Tag{
			Created:  time.Now(),
			Modified: time.Now(),
			Name:     name,
		}
		err = d.CreateTag(t)
		if err != nil {
			return nil, err
		}
	}

	d.addTagToCache(t)

	return t, nil
}

// GetTagByName returns the tag with the given name
func (d *Database) GetTagByName(name string) (*quicknote.Tag, error) {
	d.mux.Lock()
	defer d.mux.Unlock()

	if t := d.getFromTagCache(name); t != nil {
		return t, nil
	}

	sqlStr := "SELECT id, created, modified, name FROM tags WHERE name = ?;"

	stmt, err := d.db.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	t := quicknote.NewTag()
	err = stmt.QueryRow(name).Scan(&t.ID, &t.Created, &t.Modified, &t.Name)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	d.addTagToCache(t)

	return t, nil
}

// LoadNoteTags loads all the tags for the given Note
func (d *Database) LoadNoteTags(n *quicknote.Note) error {
	d.mux.Lock()
	defer d.mux.Unlock()
	return d.loadNoteTags(n)
}

func (d *Database) loadNoteTags(n *quicknote.Note) error {
	sqlStr := "SELECT id, created, modified, name FROM tags WHERE id in " +
		"(SELECT tag_id FROM note_tag WHERE note_id = ?);"

	stmt, err := d.db.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(n.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	n.Tags, err = d.loadTagsFromRows(rows)
	return err
}

// CreateTag saves the tag to the database
func (d *Database) CreateTag(t *quicknote.Tag) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	sqlStr := "INSERT INTO tags (created, modified, name) VALUES (?,?,?);"

	tx, stmt, err := d.getTxStmt(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(t.Created, t.Modified, t.Name)
	if err != nil {
		tx.Rollback()
		return err
	}

	if t.ID, err = res.LastInsertId(); err != nil {
		tx.Rollback()
		return err
	}

	d.addTagToCache(t)

	tx.Commit()
	return nil
}

func (d *Database) loadTagsFromRows(rows *sql.Rows) (quicknote.Tags, error) {
	tags := make(quicknote.Tags, 0)
	for rows.Next() {
		t := quicknote.NewTag()
		err := rows.Scan(&t.ID, &t.Created, &t.Modified, &t.Name)
		if err != nil {
			return nil, err
		}

		tags = append(tags, t)
		d.addTagToCache(t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (d *Database) createTagRal(n *quicknote.Note, tx *sql.Tx) error {
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

func (d *Database) createNoteTagRel(n *quicknote.Note, t *quicknote.Tag, tx *sql.Tx) error {
	sqlStr := "INSERT INTO note_tag (note_id, tag_id) VALUES (?,?);"

	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(n.ID, t.ID)
	if err != nil && !strings.Contains(err.Error(), "UNIQUE constraint") {
		return err
	}

	return nil
}

func (d *Database) createNoteBookTagRel(n *quicknote.Note, t *quicknote.Tag, tx *sql.Tx) error {
	sqlStr := "INSERT INTO note_book_tag (note_id, bk_id, tag_id) VALUES (?,?,?);"

	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(n.ID, n.Book.ID, t.ID)
	if err != nil && !strings.Contains(err.Error(), "UNIQUE constraint") {
		return err
	}

	return nil
}

func (d *Database) deleteTagRal(n *quicknote.Note, tx *sql.Tx) error {
	if err := d.deleteNoteTagsRel(n, tx); err != nil {
		return err
	}
	if err := d.deleteNoteNookTagsRel(n, tx); err != nil {
		return err
	}
	return nil
}

func (d *Database) deleteNoteTagsRel(n *quicknote.Note, tx *sql.Tx) error {
	sqlStr := "DELETE FROM note_tag WHERE note_id = ?"

	stmt, err := tx.Prepare(sqlStr)
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

func (d *Database) deleteNoteNookTagsRel(n *quicknote.Note, tx *sql.Tx) error {
	sqlStr := "DELETE FROM note_book_tag WHERE note_id = ?"

	stmt, err := tx.Prepare(sqlStr)
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

func (d *Database) addTagToCache(tag *quicknote.Tag) {
	d.tagNameCache[tag.Name] = tag
}

func (d *Database) delTagFromCache(tag *quicknote.Tag) {
	delete(d.tagNameCache, tag.Name)
}

func (d *Database) delTagFromCacheS(name string) {
	delete(d.tagNameCache, name)
}

func (d *Database) getFromTagCache(name string) *quicknote.Tag {
	if tag, found := d.tagNameCache[name]; found {
		return tag
	}
	return nil
}
