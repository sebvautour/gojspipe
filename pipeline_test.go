package gojspipe

import (
	"context"
	"testing"
)

func ExamplePipeline(t *testing.T) {
	ctx := context.Background()

	// Create scripts
	scriptOne, err := NewScript("script_one", `
	t.Log("script_one: root of script ran on pipeline init");

	function run() {
		t.Log("script_one: run for item " + JSON.stringify(item));
		subrun();
	}

	function subrun() {
		t.Log("script_one: also run during pipeline run, since it's called by run()");
	}
	`)
	if err != nil {
		t.Error(err)
		return
	}

	scriptTwo, err := NewScript("script_two", `
	t.Log("script_two: root of script ran on pipeline init");

	function run() {
		t.Log("script_two: run for item " + JSON.stringify(item));
		subrun();
		return false;
	}

	function subrun() {
		t.Log("script_two: also run during pipeline run, since it's called by run()");
	}
	`)
	if err != nil {
		t.Error(err)
		return
	}

	// Init new pipeline
	t.Log("- NewPipeline - start")
	p, err := NewPipeline(ctx, []*Script{scriptOne, scriptTwo}, DefaultPipeLineConfig, PipelineValue{Name: "t", Value: t})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("- NewPipeline - end")

	// runn the pipeline as needed
	t.Log("- Run - start")
	err = p.Run(PipelineValue{Name: "item", Value: map[string]string{"example": "example item"}})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("- Run - end")

	t.Log("- Run - start")
	err = p.Run(PipelineValue{Name: "item", Value: map[string]string{"example": "another example item"}})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("- Run - end")

}
