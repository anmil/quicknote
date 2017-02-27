[![GoDoc](https://godoc.org/github.com/anmil/quicknote?status.svg)](https://godoc.org/github.com/anmil/quicknote) [![Build Status](https://travis-ci.org/anmil/quicknote.svg?branch=master)](https://travis-ci.org/anmil/quicknote) [![Go Report Card](https://goreportcard.com/badge/github.com/anmil/quicknote)](https://goreportcard.com/report/github.com/anmil/quicknote)

## What is QuickNote?

Qnote allows you to quickly create and search tens of thousands of short notes.

* Create Books to organize Notes
* Create Notes using vim (overridden by Env: EDITOR)
* Create Notes from URL, auto generating notes using meta data from the website
* Create Tags in Note by prefix words with `#`
* Edit Notes
* Search Notes using Bleve or ElasticSearch QueryStringQuery or PhrasePrefix.
* Export Notes in colored text, csv, json.
* Delete, Merge, and Split Books
* Works out of the box with no configuration required.
* Supports Linux and Mac OSX
* Experimental CUI interface

Notes are stored in an SQLite database (support for more databases is coming). Searching is provided by Bleve (default) or Elasticsearch with some extra setup.

![qnote](https://cloud.githubusercontent.com/assets/1073151/23346064/fceeecba-fc64-11e6-9498-52038c853ddb.gif)

## Install

If you have not already done so, you need to setup [Golang](https://golang.org/). Than you just run

	go get github.com/anmil/quicknote/cmd/qnote

This will pull the library and build it into you Golang bin directory.

## Creating Books

Book allow you to keep related notes separated from each other, such as work notes vs personal notes. Unless stated otherwise, every action is preformed only on the working book. You can change the working book with the `-n` flag.

To create a new book

	qnote new book <book name>

## List all Books

	qnote ls books

## Deleting Books

You can delete books and all of the notes in the book.

	qnote rm book <book name>

If you want to remove a book but keep the notes. You can merge the book into another one. Merging Books takes all the notes from one book and moves them to another, than deletes the empty book.

	qnote merge <book to delete> <book to move notes to>

## Splitting Books

Books can be split in two ways, either from the results of a query or a list of Note IDs.

by query

	qnote split query <book name> <query>

This will preform a query search (using QueryStringQuery) on the working book and the results of the query will be moved to the `<book name>` (creating it if it does not exist). You do not need to specify the book in the query string. It is already added for you.

by ids

	qnote split ids <book name> <note ids...>

This does the same but instead of querying for the Notes, it moves the Notes for the given ids.

## Creating Notes

To open the editor (default Vim, overide with Env: EDITOR) and create a new note

	qnote new note

Notes are ran through a parser that take the first line as the title and the rest as the body. Any word starting with a `#` character is used as a tag for the note.

Creating a note with the following text

	This is a test #note
	notes are #cool and #fun
	one #note is never enough

will create a note with a title of `This is a test #note`, Tags `note, cool, fun`, and a body of

	notes are #cool and #fun
	one #note is never enough

You can also create a note from a URL.

	qnote new url <url>

qnote will preform a GET request on the URL. It will parse the returned HTML for the web page's `title`, `meta[name=keywords]`, and `meta[name=description]` tags. Title plus the URL is used as the title, keywords are used for the tags, and description for the body. It will open the editor with this information filled out and allow you to make changes before saving it.

## Listing Notes

To list all notes in a book

	qnote ls notes

To list all notes in all books

	qnote ls notes all

## Edit Note

To open the editor and edit a note

	qnote edit note <note id>

## Delete Note

To delete a note

	qnote rm note `<note id>`

## Searching Notes

qnote uses [Bleve](https://github.com/blevesearch/bleve) by default, but also supports  [ElasticSearch](https://www.elastic.co/), to index notes and allow for searching. ElasticSearch is recommend if you don't mind a little extra setup as it is much more powerful and faster. If you install Elasticsearch, you can edit qnote's config file located in `$HOME/.config/quicknote` on Linux and `$HOME/Library/Application Support/quicknote` on Max OSX.

To search your notes in the working Book using Phrase Prefix

	qnote search query

You can also use the more powerful query syntax QueryStringQuery ([Bleve](http://www.blevesearch.com/docs/Query-String-Query/) [ElasticSearch](https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html)) by setting the `-q` flag. Note when using the `-q` flag, the query runs on all Books.

The fields you can search on are `id`, `created`, `modifed`, `title`, `tags` (common separated list), `body`, and `book`. So, if you want to search for any notes in book "Work" that has the tag "projectx". You would run

For Bleve

	qnote search -q "+book:Work +tag:projectx"

and ElasticSearch

	qnote search -q "book:Work AND tags:projectx"

### Re-Indexing

When you create, edit, and delete notes, qnote will take care of updating the index. But, if you need to re-index for reasons such as, changing indexing providers, re-installed ElasticSearch, copying the notes database from another system. You can run

	qnote search reindex

and qnote will re-index all of the notes.


## Backing up and Restoring

The simplest way to back up your notes is to copy the qnote.db in the data directory `$HOME/.config/quicknote` on Linux and `$HOME/Library/Application Support/quicknote` on Max OSX. You can also export the list in csv or json format with the `-f` flag in the `qnote ls notes` command.

Currently, there is no way to restore notes from csv or json.

## Command Docs

All commands and flags are documents in the `help` command. Simple run `qnote help <command>` to view the description and flags for any command

Example 1:

	$ qnote help
	Qnote allows you to quickly create and search tens of thousands of short notes.

	Create Books to organize collections of notes.
	Add tags to notes for more accurate searching.
	Export your notes in text, csv, and json.

	Notes are stored in an SQLite database (support for more databases is coming).
	Searching is provided by Bleve by default, or Elasticsearch with some extra setup.

	Usage:
	  qnote [command]

	Available Commands:
	  delete      Delete Note or Notebook
	  edit        Edit Note or Notebook
	  get         Get/List Notes, Notebook, Tags
	  merge       Merge all notes from <book_name 1> into <book_name 2>
	  new         Create new Note or Notebook
	  search      Search notes
	  split       Split Book
	  version     Print the version of qnote

	Flags:
	  -n, --notebook string   Working Notebook (default "General")

	Use "qnote [command] --help" for more information about a command.

Example 2:

	$ qnote help edit book
	Edit the working Book's name. This requires re-index the Book

	Usage:
	  qnote edit book <new book_name> [flags]

	Global Flags:

Example 3:

	$ qnote get -h
	Get/List Notes, Notebook, Tags

	Usage:
	  qnote get [command]

	Aliases:
	  get, list, ls


	Available Commands:
	  book        List all books
	  note        List all notes in the working Book, or all notes for the given [note id...]
	  tag         lists all tags for the working Book

	Flags:
	  -f, --display-format string   Format to display notes in [ids, text, short, json, csv] (default "text")
	  -d, --display-order string    The order to display Notebook, Notes, Tags [asc, desc] (default "asc")
	  -s, --sort-by string          Sort notes by [id, created, modified, title] (default "modified")

	Global Flags:
	  -n, --notebook string   Working Notebook (default "General")

	Use "qnote get [command] --help" for more information about a command.
