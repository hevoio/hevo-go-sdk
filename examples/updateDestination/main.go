// Command updateDestination shows how to update destination with hevo-go.
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

	response, err := client.UpdateDestination().DestinationID(7).Config(map[string]any{"warehouse": "LOAD_WH_2"}).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", response)
}
