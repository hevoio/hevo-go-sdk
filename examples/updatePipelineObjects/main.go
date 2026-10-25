// Command updatePipelineObjects shows how to update pipeline objects with hevo-go.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hevoio/hevo-go"
)

func main() {
	client := hevo.New(os.Getenv("HEVO_ACCESS_KEY"), os.Getenv("HEVO_SECRET_KEY"))
	if err := client.Region(os.Getenv("HEVO_REGION")); err != nil {
		log.Fatal(err)
	}

	hevo.Debug(true)

	response, err := client.UpdatePipelineObjects().
		PipelineID(42).
		SelectionMode("EXCLUDE_ALL").
		Object(hevo.NewObjectSelection().Name("contacts").Selected(true)).
		Object(hevo.NewObjectSelection().Name("deals").Selected(true).LoadMode("APPEND")).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", response)
}
