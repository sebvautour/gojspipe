package gojspipe

import (
	"context"
	"errors"
	"log"
	"time"
)

// Pipeline holds all the scripts for the processing pipeline
type Pipeline struct {
	scripts []*Script
	conf    PipeLineConfig
}

// NewPipeline returns a new Pipeline, which each script given initalized
func NewPipeline(ctx context.Context, scripts []*Script, config PipeLineConfig, initialValues ...PipelineValue) (p *Pipeline, err error) {
	p = &Pipeline{
		conf:    config,
		scripts: scripts,
	}

	return p, p.initScripts(ctx, initialValues...)
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

func (p *Pipeline) initScripts(ctx context.Context, initialValues ...PipelineValue) error {
	for _, s := range p.scripts {
		if s.init {
			continue
		}

		s.VM.Interrupt = make(chan func(), 1)

		for _, v := range initialValues {
			err := s.VM.Set(v.Name, v.Value)
			log.Println("set " + s.Name + " value " + v.Name)
			if err != nil {
				return errors.New("set VM value " + v.Name + ": " + err.Error())
			}
		}

		_, err := s.runScript(ctx, p.conf.ScriptTimeout, s.script)
		if err != nil {
			return err
		}
		// don't need the raw script once it's run in the VM
		s.script = nil

		s.init = true
	}
	return nil
}

// Run executes a run() function for each script in the pipeline
// values can be given that will be added to the Otto VM
func (p *Pipeline) Run(ctx context.Context, values ...PipelineValue) (err error) {
	return p.runScripts(ctx, `run()`, values...)
}

func (p *Pipeline) runScripts(ctx context.Context, src interface{}, values ...PipelineValue) (err error) {
	for _, s := range p.scripts {
		// set values
		for _, v := range values {
			err = s.VM.Set(v.Name, v.Value)
			if err != nil {
				return errors.New("set: " + err.Error())
			}
		}

		stop, err := s.runScript(ctx, p.conf.ScriptTimeout, src)
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

// PipelineValue can be used to pass additional values accessible in the Pipeline vm
// see https://pkg.go.dev/github.com/robertkrimen/otto?tab=doc#Otto.Set
type PipelineValue struct {
	Name  string
	Value interface{}
}
