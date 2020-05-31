package gojspipe

import (
	"errors"

	"github.com/robertkrimen/otto"
)

// Script contains the name and Base64 encoded script
type Script struct {
	Name   string
	Script *otto.Script
	VM     *otto.Otto
}

// NewScript returns a new Script for the given source
// filename must always be given (for logging purposes), even if src is given too
// if src is nil, filename will be used to as the script src
// src may be a string, a byte slice, a bytes.Buffer, or an io.Reader, but it MUST always be in UTF-8.
//
// The script should contain a run() function which is executed every pipeline run
// the whole script file is executed on pipeline init, so script init logic can be added to the root of the script
func NewScript(filename string, src interface{}) (script *Script, err error) {
	if filename == "" {
		return nil, errors.New("empty string filename given")
	}

	script = &Script{
		Name: filename,
		VM:   otto.New(),
	}

	if src != nil {
		filename = ""
	}

	script.Script, err = otto.New().Compile(filename, src)
	if err != nil {
		return nil, err
	}

	return script, nil
}
