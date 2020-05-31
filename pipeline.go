package gojspipe

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

// Pipeline holds all the scripts for the processing pipeline
type Pipeline struct {
	ctx     context.Context
	scripts []*Script
	conf    PipeLineConfig
}

// NewPipeline returns a new Pipeline, which each script given initalized
func NewPipeline(ctx context.Context, scripts []*Script, config PipeLineConfig, initialValues ...PipelineValue) (p *Pipeline, err error) {
	p = &Pipeline{
		ctx:     ctx,
		conf:    config,
		scripts: scripts,
	}

	return p, p.initScripts(initialValues...)
}

// PipeLineConfig is the config params given to NewPipeline
type PipeLineConfig struct {
	// ScriptTimeout for each script executed
	ScriptTimeout time.Duration
	// ContinueOnError when true will continue processing following scripts if one of the scripts fail
	ContinueOnError bool
}

// DefaultPipeLineConfig returns the default values for PipeLineConfig
var DefaultPipeLineConfig = PipeLineConfig{
	ScriptTimeout:   time.Minute,
	ContinueOnError: false,
}

// ErrExecTimeout is used when a script exceeds the PipeLineConfig.ScriptTimeout param
var ErrExecTimeout = errors.New("Script timeout exceeded")

func (p *Pipeline) initScripts(initialValues ...PipelineValue) error {
	for _, s := range p.scripts {
		s.VM.Interrupt = make(chan func(), 1)

		for _, v := range initialValues {
			err := s.VM.Set(v.Name, v.Value)
			log.Println("set " + s.Name + " value " + v.Name)
			if err != nil {
				return errors.New("set VM value " + v.Name + ": " + err.Error())
			}
		}

		_, err := p.runScript(s.Script, s)
		if err != nil {
			return err
		}
	}
	return nil
}

// Run executes a run() function for each script in the pipeline
// values can be given that will be added to the Otto VM
func (p *Pipeline) Run(values ...PipelineValue) (err error) {
	return p.runScripts(`run()`, values...)
}

func (p *Pipeline) runScripts(src interface{}, values ...PipelineValue) (err error) {
	for _, s := range p.scripts {
		// set values
		for _, v := range values {
			err = s.VM.Set(v.Name, v.Value)
			if err != nil {
				return errors.New("set: " + err.Error())
			}
		}

		stop, err := p.runScript(src, s)
		if err != nil && p.conf.ContinueOnError == false {
			return err
		}
		if err != nil && p.conf.ContinueOnError == true {
			// log
		}

		if stop {
			// log
			return nil
		}

	}
	return nil
}

func (p *Pipeline) runScript(src interface{}, s *Script) (stop bool, err error) {
	interruptCtx, cancelInterrupt := context.WithCancel(p.ctx)

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

	go p.runInterrupt(interruptCtx, s.VM.Interrupt)

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

func (p *Pipeline) runInterrupt(ctx context.Context, interruptCh chan func()) {
	select {
	case <-time.After(p.conf.ScriptTimeout):
		interruptCh <- func() {
			panic(ErrExecTimeout)
		}
	case <-ctx.Done():
	}

}

// PipelineValue can be used to pass additional values accessible in the Pipeline vm
// see https://pkg.go.dev/github.com/robertkrimen/otto?tab=doc#Otto.Set
type PipelineValue struct {
	Name  string
	Value interface{}
}
