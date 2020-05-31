package gojspipe

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/robertkrimen/otto"
)

// Script contains the name and Base64 encoded script
type Script struct {
	Name   string
	script *otto.Script
	VM     *otto.Otto
	init   bool
}

// NewScript returns a new Script for the given source
//
// filename must always be given (for logging purposes), even if src is given too, if src is nil, filename will be used to as the script src
//
// src may be a string, a byte slice, a bytes.Buffer, or an io.Reader, but it MUST always be in UTF-8.
//
// The script needs to contain a run() function which is executed every pipeline run
//
// the whole script file is executed on pipeline init, so script init logic can be added to the root of the script
func NewScript(filename string, src interface{}) (script *Script, err error) {
	if filename == "" {
		return nil, errors.New("empty string filename given")
	}

	script = &Script{
		Name: filename,
		VM:   otto.New(),
		init: false,
	}

	if src != nil {
		filename = ""
	}

	script.script, err = otto.New().Compile(filename, src)
	if err != nil {
		return nil, err
	}

	return script, nil
}

func (s *Script) runScript(ctx context.Context, timeout time.Duration, src interface{}) (stop bool, err error) {
	interruptCtx, cancelInterrupt := context.WithCancel(ctx)

	defer func() {
		cancelInterrupt()
		if caught := recover(); caught != nil {
			if caught == ErrExecTimeout {
				err = ErrExecTimeout
				return
			}
			panic(caught) // Something else happened, repanic!
		}

	}()

	go s.runInterrupt(interruptCtx, timeout, s.VM.Interrupt)

	v, err := s.VM.Run(src)
	if err != nil {
		vmctx := s.VM.Context()
		return false, fmt.Errorf("[%v:%v:%v] %v", s.Name, vmctx.Line, vmctx.Column, err.Error())
	}

	if v.IsBoolean() {
		vb, err := v.ToBoolean()
		if err != nil {
			return false, err
		}
		return vb, nil
	}

	return false, nil
}

func (s *Script) runInterrupt(ctx context.Context, timeout time.Duration, interruptCh chan func()) {
	select {
	case <-time.After(timeout):
		interruptCh <- func() {
			panic(ErrExecTimeout)
		}
	case <-ctx.Done():
	}

}
