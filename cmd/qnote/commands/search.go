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

package commands

import (
	"fmt"
	"strings"

	"github.com/anmil/quicknote/cmd/shared/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Command line variables
var (
	resultsLimit     int
	resultsOffset    int
	queryStringQuery bool
)

func init() {
	RootCmd.AddCommand(SearchCmd)
	SearchCmd.AddCommand(SearchReindexCmd)

	viper.SetDefault("search_results_limit", "15")
	viper.SetDefault("raw_query", "false")

	SearchCmd.PersistentFlags().IntVarP(&resultsLimit, "limit", "l",
		viper.GetInt("search_results_limit"), "Number of results to return")

	SearchCmd.PersistentFlags().IntVarP(&resultsOffset, "offset", "o", 0, "Start point is the result, use for paging")
	SearchCmd.PersistentFlags().BoolVarP(&queryStringQuery, "query-string-query", "q", viper.GetBool("query_string_query"),
		"By default qnote will alter the query to include the working notebook tag. Set this to to disable action.")

	SearchCmd.PersistentFlags().StringVarP(&displayFormat, "display-format", "f", viper.GetString("display_format"),
		fmt.Sprintf("Format to display notes in [%s]", strings.Join(displayFormatOptions, ", ")))

	SearchCmd.PersistentFlags().BoolVarP(&displayTextOneResult, "text-single-result", "", viper.GetBool("display_text_for_one_result"),
		fmt.Sprintf("Display in text mode when there is only one result"))
}

// SearchCmd Search notes
var SearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search notes",
	Long: `Search all notes in the working Book (see '-n').

Query syntax depends on the index provider that is configured. This command
(when '-q' is not given) uses a Phrase Prefix query. Results match on all
given words in the query string with the last word used as a prefix. For
better documentation on how this works. See the index providers docs

Bleve (default): http://www.blevesearch.com/docs/Query/
ElasticSearch: https://www.elastic.co/guide/en/elasticsearch/guide/current/_query_time_search_as_you_type.html

To use the QueryStringQuery syntax set the '-q' flag
	Example (ElasticSearch): title:term1 AND tags:term2 NOT (body:term3 OR body:term4)
`,
	Run: searchCmdRun,
}

func searchCmdRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		exitValidationError("invalid query string", cmd)
	}
	query := args[0]

	var ids []int64
	var total uint64
	var err error

	if queryStringQuery {
		ids, total, err = idxConn.SearchNote(query, resultsLimit, resultsOffset)
	} else {
		ids, total, err = idxConn.SearchNotePhrase(query, workingNotebook, "asc", resultsLimit, resultsOffset)
	}
	exitOnError(err)

	notes, err := dbConn.GetAllNotesByIDs(ids)
	exitOnError(err)

	if displayFormat == "short" && displayTextOneResult && len(notes) == 1 {
		displayFormat = "text"
	}

	err = utils.PrintNotes(notes, displayFormat)
	exitOnError(err)
	fmt.Printf("\nShowing %d-%d of %d\n", resultsOffset, resultsOffset+len(ids), total)
}

// SearchReindexCmd Re-indexes all Notes in all Books
var SearchReindexCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Re-indexes all Notes in all Books",
	Long: `Use this command to re-index all of your Notes from all Books

If you change the index provider, copy or replace the qnote.db. You need to
call this in order to use the search command.

Re-indexing can take several minutes depending on the number of notes and
the index provider used.`,
	Run: searchReindexCmdRun,
}

func searchReindexCmdRun(cmd *cobra.Command, args []string) {
	notes, err := dbConn.GetAllNotes(sortBy, displayOrder)
	exitOnError(err)

	err = idxConn.IndexNotes(notes)
	exitOnError(err)

	fmt.Printf("Finished indexing notes (%d)\n", len(notes))
}
