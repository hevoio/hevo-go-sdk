// Command listPipelines shows how to list pipelines with hevo-go.
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

	response, err := client.ListPipelines().States("Initialized").Limit(50).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", response)
}
