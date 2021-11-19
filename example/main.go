package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/probakowski/go-airly"
	"log"
	"time"
)

func main() {
	lon := flag.Float64("lon", 0, "Longitude")
	lat := flag.Float64("lat", 0, "Latitude")
	key := flag.String("key", "", "API key")
	installation := flag.Int("installation", -1, "Installation ID to get measurements from, -1 means longitude and latitude will be used")
	language := flag.String("lang", "en", "Language, en or pl")
	cloudId := flag.String("cloudId", "", "Elasticsearch Cloud ID")
	user := flag.String("user", "", "Elasticsearch user")
	password := flag.String("password", "", "Elasticsearch password")
	flag.Parse()

	cfg := elasticsearch.Config{
		CloudID:  *cloudId,
		Username: *user,
		Password: *password,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	air := airly.Client{
		Key:      *key,
		Language: *language,
	}

	for {
		var measurements airly.Measurements
		if *installation == -1 {
			measurements, err = air.NearestMeasurements(airly.Location{Latitude: *lat, Longitude: *lon})
		} else {
			measurements, err = air.InstallationMeasurements(*installation)
		}
		if err != nil {
			log.Printf("Error getting measurements %s\n", err)
			time.Sleep(15 * time.Minute)
			continue
		}

		data, err := json.Marshal(measurements)
		if err != nil {
			log.Printf("Error serializing measurements %s\n", err)
			time.Sleep(15 * time.Minute)
			continue
		}

		req := esapi.IndexRequest{
			Index:   "airly",
			Body:    bytes.NewReader(data),
			Refresh: "true",
		}

		res, err := req.Do(context.Background(), es)
		if err != nil {
			log.Printf("Error getting response: %s\n", err)
			time.Sleep(15 * time.Minute)
			continue
		}
		if res.IsError() {
			log.Printf("[%s] Error indexing document", res.Status())
		} else {
			// Deserialize the response into a map.
			var r map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				log.Printf("Error parsing the response body: %s", err)
			} else {
				// Print the response status and indexed document version.
				log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
			}
		}
		_ = res.Body.Close()

		time.Sleep(15 * time.Minute)
	}
}
