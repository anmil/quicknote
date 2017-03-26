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

package encoding

// Binary file exporter and importer
//
// A QNOT file consists of a Header and one or more records. Every QNOT file must
// start with the Header and can have any number of records following it. There
// are no more records once the EOF is reached.
//
// All multi-byte integer fields are encoded in Big-endian order.
//
// All strings are pre-pended with a 8 byte (uint64) number containing the string
// length in bytes.
//
// Base format looks as follows
//
// 	+=========================================+
// 	|                Header                   |
// 	+=========================================+
// 	+=========================================+
// 	|                Record                   |
// 	+=========================================+
// 	+=========================================+
// 	|                Record                   |
// 	+=========================================+
// 	+=========================================+
// 	|                 ...                     |
// 	+=========================================+
//
// The Header block starts with the magic string "QNOT", followed by the format
// version, and a timestamp of when the file was created.
//
// 	| 4 byte magic "QNOT" | 4 byte version (uint32) | 8 byte timestamp (uint64) |
//
// There are three types of records; Book, Tag, and Note. The first byte of a
// record specifies the record type.
//
// All Book and Tag records that a Note record references MUST appear before
// the Note record that referenced it.
//
// Record Types
// Book = 0 (0x00)
// Tag  = 1 (0x01)
// Note = 2 (0x02)
//
// Book Record
//
// 	| 1 byte record type "0"                  |
// 	| 8 byte book ID (uint64)                 |
// 	| 8 byte book Created timestamp (uint64)  |
// 	| 8 byte book Modified timestamp (uint64) |
// 	| 8 byte Name string length (uint64)      |
// 	| varlen book Name byte string            |
//
// Tag Record
//
// 	| 1 byte record type "1"                 |
// 	| 8 byte tag ID (uint64)                 |
// 	| 8 byte tag Created timestamp (uint64)  |
// 	| 8 byte tag Modified timestamp (uint64) |
// 	| 8 byte Name string length (uint64)     |
// 	| varlen tag Name byte string            |
//
// Note Record
//
// 	| 1 byte record type "2"                 |
// 	| 8 byte tag ID (uint64)                 |
// 	| 8 byte tag Created timestamp (uint64)  |
// 	| 8 byte tag Modified timestamp (uint64) |
// 	| 8 byte Type string length (uint64)     |
// 	| varlen tag Type byte string            |
// 	| 8 byte Title string length (uint64)    |
// 	| varlen tag Title byte string           |
// 	| 8 byte Body string length (uint64)     |
// 	| varlen tag Body byte string            |
// 	| 8 byte book ID (uint64)                |
// 	| 8 byte number of tags (uint64)         |
// 	| 8 byte tag ID (uint64)                 | <- repeats for each tag
//

import (
	"bytes"
	"errors"
	"io"
	"time"

	"github.com/anmil/quicknote/note"

	"encoding/binary"
)

// ErrHeaderNotWritten indicates that the Header was not written before attempting
// to write a record.
var ErrHeaderNotWritten = errors.New("Header have not been written yet")

// ErrHeaderNotParsed indicates that the Header was not read before attempting
// to read a record
var ErrHeaderNotParsed = errors.New("Header have not been parsed yet")

// ErrInvalidRecordType indicates a invalid record was encountered. See RecordType
// for a list of valid records types.
var ErrInvalidRecordType = errors.New("Encountered an invalid record type")

// ErrInvalidString indicates the parser encountered a string that does not
// meets the requirements of Entity's field (such as a string was to long)
var ErrInvalidString = errors.New("Encountered an invalid record string")

// ErrInvalidOrCorruptedBinaryFormat indicates that parser encountered an problem
// with the QNOT stream and is unable to continue.
var ErrInvalidOrCorruptedBinaryFormat = errors.New("Invalid or corrupted binary")

// ErrBookNoteFound indicates that the parser encountered a Note record that
// referenced a Book the parse has not parsed yet.
var ErrBookNoteFound = errors.New("Note how unknown book")

// ErrTagNoteFound indicates that the parser encountered a Note record that
// referenced a Tag the parse has not parsed yet.
var ErrTagNoteFound = errors.New("Note how unknown tag")

// MagicStr binary magic string
var MagicStr = "QNOT"

// CurrentVersion the current format version used for encoding and decoding
var CurrentVersion uint32 = 1

// HeaderLen length of the header block
var HeaderLen = 16

// RecordType byte indicating the type of record
type RecordType byte

// List of record types
var (
	Book RecordType = 0
	Tag  RecordType = 1
	Note RecordType = 2
)

// BinaryEncoder encodes a Note into the QNOT format
type BinaryEncoder struct {
	w io.Writer

	headerWritten bool

	wBooks map[int64]bool
	wTags  map[int64]bool
}

// NewBinaryEncoder returns a new BinaryEncoder
func NewBinaryEncoder(w io.Writer) *BinaryEncoder {
	return &BinaryEncoder{
		w:      w,
		wBooks: make(map[int64]bool),
		wTags:  make(map[int64]bool),
	}
}

// WriteHeader writes the Header block to w.
// This must be called before any Note can be encoded
// Calling WriteHeader() more than once has no effect.
func (b *BinaryEncoder) WriteHeader() (uint64, error) {
	// If the header has already been written, we simple just return
	if b.headerWritten {
		return 0, nil
	}

	buff := &bytes.Buffer{}

	if _, err := buff.Write([]byte(MagicStr)); err != nil {
		return 0, err
	}
	if err := writeInt32(buff, int32(CurrentVersion)); err != nil {
		return 0, err
	}
	if err := writeTime(buff, time.Now()); err != nil {
		return 0, err
	}

	data := buff.Bytes()
	bw, err := b.w.Write(data)

	b.headerWritten = true
	return uint64(bw), err
}

// WriteNote encodes a Note, it's Book, and Tags then writes them to w.
// Books and Tags are only encoded and written to w once. Meaning, if
// they are encountered again in a different note, they will not be
// encoded again.
func (b *BinaryEncoder) WriteNote(n *note.Note) (uint64, error) {
	if !b.headerWritten {
		return 0, ErrHeaderNotWritten
	}

	var bytesWritten uint64

	bw, err := b.writeBook(n.Book)
	if err != nil {
		return bytesWritten, err
	}
	bytesWritten += bw

	for _, tag := range n.Tags {
		bw, err = b.writeTag(tag)
		if err != nil {
			return bytesWritten, err
		}
		bytesWritten += bw
	}

	bw, err = b.writeNote(n)
	if err != nil {
		return bytesWritten, err
	}
	bytesWritten += bw

	return bytesWritten, nil
}

func (b *BinaryEncoder) writeBook(bk *note.Book) (uint64, error) {
	// Check if we have already written this book to the stream
	if _, found := b.wBooks[bk.ID]; found {
		return 0, nil
	}

	buff := &bytes.Buffer{}

	if _, err := buff.Write([]byte{byte(Book)}); err != nil {
		return 0, err
	}
	if err := writeInt64(buff, bk.ID); err != nil {
		return 0, err
	}
	if err := writeTime(buff, bk.Created); err != nil {
		return 0, err
	}
	if err := writeTime(buff, bk.Modified); err != nil {
		return 0, err
	}
	if err := writeString(buff, bk.Name); err != nil {
		return 0, err
	}
	wb, err := b.writeBuffer(buff)

	b.wBooks[bk.ID] = true
	return wb, err
}

func (b *BinaryEncoder) writeTag(tg *note.Tag) (uint64, error) {
	// Check if we have already written this book to the stream
	if _, found := b.wTags[tg.ID]; found {
		return 0, nil
	}

	buff := &bytes.Buffer{}

	if _, err := buff.Write([]byte{byte(Tag)}); err != nil {
		return 0, err
	}
	if err := writeInt64(buff, tg.ID); err != nil {
		return 0, err
	}
	if err := writeTime(buff, tg.Created); err != nil {
		return 0, err
	}
	if err := writeTime(buff, tg.Modified); err != nil {
		return 0, err
	}
	if err := writeString(buff, tg.Name); err != nil {
		return 0, err
	}
	wb, err := b.writeBuffer(buff)

	b.wTags[tg.ID] = true
	return wb, err
}

func (b *BinaryEncoder) writeNote(n *note.Note) (uint64, error) {
	buff := &bytes.Buffer{}

	if _, err := buff.Write([]byte{byte(Note)}); err != nil {
		return 0, err
	}
	if err := writeInt64(buff, n.ID); err != nil {
		return 0, err
	}
	if err := writeTime(buff, n.Created); err != nil {
		return 0, err
	}
	if err := writeTime(buff, n.Modified); err != nil {
		return 0, err
	}
	if err := writeString(buff, n.Type); err != nil {
		return 0, err
	}
	if err := writeString(buff, n.Title); err != nil {
		return 0, err
	}
	if err := writeString(buff, n.Body); err != nil {
		return 0, err
	}
	if err := writeInt64(buff, n.Book.ID); err != nil {
		return 0, err
	}
	if err := writeInt64Slice(buff, n.GetTagIDsArray()); err != nil {
		return 0, err
	}

	return b.writeBuffer(buff)
}

func writeString(buff io.Writer, s string) error {
	data := []byte(s)

	dataLen := make([]byte, 8)
	binary.BigEndian.PutUint64(dataLen, uint64(len(data)))

	if _, err := buff.Write(dataLen); err != nil {
		return err
	}
	_, err := buff.Write(data)
	return err
}

func writeTime(buff io.Writer, t time.Time) error {
	return writeInt64(buff, t.Unix())
}

func writeInt32(buff io.Writer, i int32) error {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(i))
	_, err := buff.Write(data)
	return err
}

func writeInt64Slice(buff io.Writer, is []int64) error {
	l := int64(len(is))
	if err := writeInt64(buff, l); err != nil {
		return err
	}

	for _, i := range is {
		if err := writeInt64(buff, i); err != nil {
			return err
		}
	}

	return nil
}

func writeInt64(buff io.Writer, i int64) error {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(i))
	_, err := buff.Write(data)
	return err
}

func (b *BinaryEncoder) writeBuffer(buff *bytes.Buffer) (uint64, error) {
	var bytesWritten uint64
	data := buff.Bytes()

	bw, err := b.w.Write(data)
	bytesWritten += uint64(bw)

	return bytesWritten, err
}

// BinaryDecoder Decodes a QNOT file into a list of Notes
type BinaryDecoder struct {
	r io.Reader

	Header *Header
	Err    error

	wBooks map[int64]*note.Book
	wTags  map[int64]*note.Tag
}

// Header QNOT file header block
type Header struct {
	Version int32
	Created time.Time
}

// NewBinaryDecoder returns a new BinaryDecoder
func NewBinaryDecoder(r io.Reader) *BinaryDecoder {
	return &BinaryDecoder{
		r:      r,
		wBooks: make(map[int64]*note.Book),
		wTags:  make(map[int64]*note.Tag),
	}
}

// ParseHeader parses the Header block
// This must be called before calling ParseNotes()
// It is an error not to
func (d *BinaryDecoder) ParseHeader() error {
	buff := make([]byte, 4)
	_, err := io.ReadFull(d.r, buff)
	if err != nil {
		return err
	}

	magic := string(buff)
	if magic != MagicStr {
		return ErrInvalidOrCorruptedBinaryFormat
	}

	version, err := readInt32(d.r)
	if err != nil {
		return err
	}

	created, err := readTime(d.r)
	if err != nil {
		return err
	}

	d.Header = &Header{
		Version: version,
		Created: created,
	}

	return nil
}

// ParseNotes starts the parser and returns a Note channel
// Must call ParseHeader() first or an error is returned
func (d *BinaryDecoder) ParseNotes() (<-chan *note.Note, error) {
	if d.Header == nil {
		return nil, ErrHeaderNotParsed
	}

	notes := make(chan *note.Note, 1024)
	go d.parseNotes(notes)
	return notes, nil
}

func (d *BinaryDecoder) parseNotes(out chan *note.Note) {
	for {
		t, err := readRecordType(d.r)
		if err == io.EOF {
			break
		} else if err != nil {
			d.Err = err
			break
		}

		switch t {
		case Note:
			n, err := d.parseNote()
			if err != nil {
				d.Err = err
				break
			}
			out <- n
		case Book:
			bk, err := d.parseBook()
			if err != nil {
				d.Err = err
				break
			}
			d.wBooks[bk.ID] = bk
		case Tag:
			tag, err := d.parseTag()
			if err != nil {
				d.Err = err
				break
			}
			d.wTags[tag.ID] = tag
		default:
			d.Err = ErrInvalidRecordType
			break
		}
	}

	close(out)
}

func (d *BinaryDecoder) parseBook() (*note.Book, error) {
	var err error
	bk := note.NewBook()

	if bk.ID, err = readInt64(d.r); err != nil {
		return nil, err
	}
	if bk.Created, err = readTime(d.r); err != nil {
		return nil, err
	}
	if bk.Modified, err = readTime(d.r); err != nil {
		return nil, err
	}
	if bk.Name, err = readString(d.r); err != nil {
		return nil, err
	}

	return bk, nil
}

func (d *BinaryDecoder) parseTag() (*note.Tag, error) {
	var err error
	t := note.NewTag()

	if t.ID, err = readInt64(d.r); err != nil {
		return nil, err
	}
	if t.Created, err = readTime(d.r); err != nil {
		return nil, err
	}
	if t.Modified, err = readTime(d.r); err != nil {
		return nil, err
	}
	if t.Name, err = readString(d.r); err != nil {
		return nil, err
	}

	return t, nil
}

func (d *BinaryDecoder) parseNote() (*note.Note, error) {
	var err error
	n := note.NewNote()

	if n.ID, err = readInt64(d.r); err != nil {
		return nil, err
	}
	if n.Created, err = readTime(d.r); err != nil {
		return nil, err
	}
	if n.Modified, err = readTime(d.r); err != nil {
		return nil, err
	}
	if n.Type, err = readString(d.r); err != nil {
		return nil, err
	}
	if n.Title, err = readString(d.r); err != nil {
		return nil, err
	}
	if n.Body, err = readString(d.r); err != nil {
		return nil, err
	}

	bkID, err := readInt64(d.r)
	if err != nil {
		return nil, err
	}

	bk, found := d.wBooks[bkID]
	if !found {
		return nil, ErrBookNoteFound
	}
	n.Book = bk

	tagIDs, err := readInt64Slice(d.r)
	if err != nil {
		return nil, err
	}

	for _, tid := range tagIDs {
		tag, found := d.wTags[tid]
		if !found {
			return nil, ErrTagNoteFound
		}
		n.Tags = append(n.Tags, tag)
	}

	return n, nil
}

func readRecordType(rd io.Reader) (RecordType, error) {
	buff := make([]byte, 1)
	_, err := io.ReadFull(rd, buff)
	if err != nil {
		return 0, err
	}
	return RecordType(buff[0]), nil
}

func readString(rd io.Reader) (string, error) {
	strLen, err := readInt64(rd)
	if err != nil {
		return "", err
	} else if strLen > note.MaxStringLen {
		return "", ErrInvalidString
	}
	buff := make([]byte, strLen)
	_, err = io.ReadFull(rd, buff)
	if err != nil {
		return "", err
	}
	return string(buff), nil
}

func readTime(rd io.Reader) (time.Time, error) {
	ts, err := readInt64(rd)
	if err != nil {
		return time.Time{}, err
	}
	t := time.Unix(ts, 0)
	return t, nil
}

func readInt32(rd io.Reader) (int32, error) {
	buff := make([]byte, 4)
	_, err := io.ReadFull(rd, buff)
	if err != nil {
		return 0, err
	}
	i := binary.BigEndian.Uint32(buff)
	return int32(i), nil
}

func readInt64Slice(rd io.Reader) ([]int64, error) {
	l, err := readInt64(rd)
	if err != nil {
		return nil, err
	}

	is := make([]int64, l)
	for i := int64(0); i < l; i++ {
		c, err := readInt64(rd)
		if err != nil {
			return nil, err
		}
		is[i] = c
	}

	return is, nil
}

func readInt64(rd io.Reader) (int64, error) {
	buff := make([]byte, 8)
	_, err := io.ReadFull(rd, buff)
	if err != nil {
		return 0, err
	}
	i := binary.BigEndian.Uint64(buff)
	return int64(i), nil
}
