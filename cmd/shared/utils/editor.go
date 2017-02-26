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

package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const defaultEditor = "/usr/bin/vim"
const editorEnvVariable = "EDITOR"

// Editor provides an interface to open text editors
// and retrieve the text entered
type Editor struct {
	tmpFile *os.File
	text    string
	open    bool
}

// NewEditor returns a new Editor
func NewEditor() (*Editor, error) {
	tmpfile, err := ioutil.TempFile("", "rt")
	if err != nil {
		return nil, nil
	}

	return &Editor{tmpFile: tmpfile, open: false}, nil
}

// SetText sets the text to be opened in the Editor.
// This has no effect after the editor is opened
func (e *Editor) SetText(text string) {
	if !e.open {
		e.tmpFile.Write([]byte(text))
	}
}

// Open opens a text editor specified by defaultEditor (overridden by editorEnvVariable)
// Than waits for it to close and saves the text that was entered.
func (e *Editor) Open() error {
	e.open = true

	editorCmd := os.Getenv(editorEnvVariable)
	if editorCmd == "" {
		editorCmd = defaultEditor
	}

	cmdArgs := make([]string, 0)
	parts := strings.Split(editorCmd, " ")
	if len(parts) == 0 {
		e.open = false
		return errors.New("Editor command found")
	} else if len(parts) > 1 {
		cmdArgs = append(cmdArgs, parts[1:]...)
	}
	cmdP := parts[0]
	cmdArgs = append(cmdArgs, e.tmpFile.Name())

	cmd := exec.Command(cmdP, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		e.open = false
		return nil
	}

	f, err := os.Open(e.tmpFile.Name())
	if err != nil {
		e.open = false
		return nil
	}
	defer f.Close()

	body, err := ioutil.ReadAll(f)
	e.text = string(body)

	e.open = false
	return nil
}

// Text returns the text the user entered in the editor
func (e *Editor) Text() string {
	return e.text
}

// Close closes the temp file and removes it
func (e *Editor) Close() {
	e.tmpFile.Close()
	os.Remove(e.tmpFile.Name())
}
