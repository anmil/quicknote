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

package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/anmil/quicknote/cmd/shared/utils"
	"github.com/anmil/quicknote/db"
	"github.com/anmil/quicknote/index"
	"github.com/anmil/quicknote/note"
	"github.com/spf13/viper"
)

// TestingMode TODO
var TestingMode = false

// DataDirectory TODO
var DataDirectory = utils.GetDataDirectory()

// IndexProvider TODO
var IndexProvider string

func init() {
	isTestingMode()

	// Make sure our data directory exists
	err := utils.EnsureDirectoryExists(DataDirectory)
	if err != nil {
		log.Fatalln(err)
	}

	// Created the config file if it does not exists
	configFilePath := path.Join(DataDirectory, "qnote.yaml")
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		bt := []byte(defaultConfigFileText)
		err := ioutil.WriteFile(configFilePath, bt, 0600)
		if err != nil {
			log.Fatalln(err)
		}
	}

	viper.SetConfigName("qnote")
	viper.AddConfigPath(DataDirectory)
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}

	viper.SetDefault("default_notebook", "General")
	viper.SetDefault("db_provider", "sqlite")

	viper.SetDefault("index_provider", "bleve")
	viper.SetDefault("bleve_shard_count", "16")
	viper.SetDefault("elastic_url", "http://127.0.0.1:9200")
	viper.SetDefault("elastic_index_name", "qnote")

	IndexProvider = viper.GetString("index_provider")
}

func isTestingMode() {
	// This will be true when running "go run main.go"
	// I don't want to mess up my own qnote's notebooks :)
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if strings.Contains(wd, "go/src") {
		DataDirectory = path.Join(os.Getenv("HOME"), ".config", "quicknote-dev")

		TestingMode = true

	}
}

// GetDBConn gets a new Database connection for the config provider
func GetDBConn() (db.DB, error) {
	switch viper.GetString("db_provider") {
	case "sqlite":
		return getSqliteDBConn()
	case "postgres":
		return getPostgresDBConn()
	default:
		return nil, errors.New("Unsupported database provider")
	}
}

func getSqliteDBConn() (db.DB, error) {
	fp := path.Join(DataDirectory, "notes.db")
	d, err := db.NewDatabase("sqlite", fp)
	if err != nil {
		return nil, err
	}
	return d, err
}

func getPostgresDBConn() (db.DB, error) {
	name := viper.GetString("postgres.name")
	host := viper.GetString("postgres.host")
	port := viper.GetString("postgres.port")
	user := viper.GetString("postgres.user")
	pass := viper.GetString("postgres.pass")
	sslmode := viper.GetString("postgres.sslmode")

	d, err := db.NewDatabase("postgres", name, host, port, user, pass, sslmode)
	if err != nil {
		return nil, err
	}
	return d, err
}

// GetIndexConn gets a new Index connection for the config provider
func GetIndexConn() (index.Index, error) {
	switch IndexProvider {
	case "bleve":
		return getBleveConn()
	case "elastic":
		return getESConn()
	default:
		return nil, errors.New("Unsupported index provider")
	}
}

func getBleveConn() (index.Index, error) {
	shareds := viper.GetString("bleve_shard_count")
	idxConn, err := index.NewIndex("bleve", DataDirectory, shareds)
	if err != nil {
		return nil, err
	}
	return idxConn, nil
}

func getESConn() (index.Index, error) {
	url := viper.GetString("elastic_url")
	indexName := viper.GetString("elastic_index_name")
	idxConn, err := index.NewIndex("elastic", url, indexName)
	if err != nil {
		return nil, err
	}
	return idxConn, nil
}

// GetWorkingBook gets the config working Book
func GetWorkingBook(db db.DB, bkName string) (*note.Book, error) {
	if bkName == viper.GetString("default_book") {
		return db.GetOrCreateBookByName(bkName)
	}
	return db.GetBookByName(bkName)
}

// Default config file for QuickNote
var defaultConfigFileText = `
# Default Book to use with call commands
default_book: General

# Order to display results when getting multiple.
# Such as sorting by created or modified dates
# This is used in conjuration with "order_by"
# Options: asc, desc
display_order: asc

# Specified what to order the results by when getting
# multiple results. See also "display_order"
# Options: id, created, modified, title
order_by: modified

# Number of results to return for search queries
search_results_limit: 15

# By default qnote will alter the search query to
# include the working notebook tag. Set this to
# true to disable action.
raw_query: false

# Database provider
# Options: sqlite
# TODO: Will be adding support for Postgres in the future
db_provider: sqlite
# db_provider: postgres

# PostgreSQl settings
# postgres:
#   name: qnote
#   user: qnote
#   pass: *****
#   host: localhost
#   port: 5432
#   sslmode: disable

# This options are not required for sqlite
# db_host: 127.0.0.1
# db_port: 5432
# db_name: qnote
# db_user: qnote
# db_pass: *****

# Indexing provider
# Currently only Bleve and Elasticsearch
# Bleve: http://www.blevesearch.com/docs/Query-String-Query/
# Elasticsearch: https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html
# are supported.
index_provider: bleve
# index_provider: elastic

# Qnote will split notes across multiple Bleve indexes
# bleve_shard_count is the number of indexes to use.
# The number of shards you should use depends on you
# read/write disc speed. If the default does not work
# you will need to experiment to find the correct value.
#
# If you change this value, you will need to delete the .bleve
# index folders in the data directory and call
# "qnote search reindex" to rebuilt the indexes
bleve_shard_count: 16

# elastic_url: http://127.0.0.1:9200
# elastic_index_name: qnote
`
