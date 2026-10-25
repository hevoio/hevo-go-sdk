// Command getPipelineJob shows how to get pipeline job with hevo-go.
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

	response, err := client.GetPipelineJob().PipelineID(42).JobID("9b2c3d4e-0000-4000-8000-000000000000").
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", response)
}
