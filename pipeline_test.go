package gojspipe

import (
	"context"
	"log"
)

func ExampleNewPipeline() {
	ctx := context.Background()

	// Create scripts
	scriptOne, err := NewScript("script_one", `
	log.Println("script_one: root of script ran on pipeline init");

	function run() {
		log.Println("script_one: run for item " + JSON.stringify(item));
		subrun();
	}

	function subrun() {
		log.Println("script_one: also run during pipeline run, since it's called by run()");
	}
	`)
	if err != nil {
		log.Println(err)
		return
	}

	scriptTwo, err := NewScript("script_two", `
	log.Println("script_two: root of script ran on pipeline init");

	function run() {
		log.Println("script_two: run for item " + JSON.stringify(item));
		subrun();
		return false;
	}

	function subrun() {
		log.Println("script_two: also run during pipeline run, since it's called by run()");
	}
	`)
	if err != nil {
		log.Println(err)
		return
	}

	// Init new pipeline
	log.Println("- NewPipeline - start")
	p, err := NewPipeline(ctx, []*Script{scriptOne, scriptTwo}, DefaultPipeLineConfig, PipelineValue{Name: "log", Value: log.Logger{}})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("- NewPipeline - end")

	// runn the pipeline as needed
	log.Println("- Run - start")
	err = p.Run(PipelineValue{Name: "item", Value: map[string]string{"example": "example item"}})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("- Run - end")

	log.Println("- Run - start")
	err = p.Run(PipelineValue{Name: "item", Value: map[string]string{"example": "another example item"}})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("- Run - end")

}
