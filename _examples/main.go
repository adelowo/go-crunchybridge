package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	gocrunchybridge "github.com/adelowo/go-crunchybridge"
)

func main() {

	client, err := gocrunchybridge.New(
		gocrunchybridge.WithAPIKey(gocrunchybridge.APIKey(os.Getenv("INTEGRATION_API_KEY"))))

	if err != nil {
		panic(err)
	}

	_, err = client.Provider.Get(context.Background(), gocrunchybridge.FetchProviderOptions{
		Provider: gocrunchybridge.ClusterProviderAws,
	})

	plan, err := client.Provider.Get(context.Background(), gocrunchybridge.FetchProviderOptions{
		Provider: gocrunchybridge.ClusterProviderGcp,
	})

	json.NewEncoder(os.Stdout).Encode(plan)

	log.Fatal(err)

	cluster, err := client.Cluster.Create(context.Background(), &gocrunchybridge.CreateClusterOptions{
		PlanID:            "hobby-0",
		TeamID:            "zar5556utjb3hkevt5dkxj2o4i",
		RegionID:          "eu-west-2",
		ProviderID:        gocrunchybridge.ClusterProviderAws,
		StorageSize:       10,
		Environment:       gocrunchybridge.ClusterEnvironmentProduction,
		HighlyAvailable:   false,
		PostgresVersionID: 16,
	})

	if err != nil {
		panic(err)
	}

	json.NewEncoder(os.Stdout).Encode(cluster)

	clusters, err := client.Cluster.List(context.Background(), gocrunchybridge.ListClusterOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println()
	json.NewEncoder(os.Stdout).Encode(clusters)

	err = client.Cluster.Delete(context.Background(), gocrunchybridge.FetchClusterOptions{
		ID: cluster.ID,
	})
	if err != nil {
		panic(err)
	}
}
