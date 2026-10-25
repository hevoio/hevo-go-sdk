// Command createPipeline shows how to create pipeline with hevo-go.
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

	response, err := client.CreatePipeline().
		Name("HubSpot to Snowflake").
		Source(hevo.NewSource().
			Name("HubSpot").
			ConnectorID("hubspot").
			Config(map[string]any{
				"authentication_method":    "PRIVATE_APP_TOKEN",
				"private_app_access_token": "pat-na1-xxxxxxxx",
			})).
		DestinationID(7).
		DestinationPrefix("hubspot_").
		ReplicationType("HISTORICAL_AND_INCREMENTAL").
		LoadMode("MERGE").
		SchemaEvolution("ALLOW_ALL").
		Schedule(hevo.NewSchedule().SyncType("SCHEDULED").FrequencyMinutes(60)).
		FailureHandlingPolicy(hevo.NewFailureHandlingPolicy().ObjectFailureLevel("PARTIAL").ObjectFailureThreshold(5)).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", response)
}
