package airly

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Installation struct {
	Id        int      `json:"id"`
	Location  Location `json:"location"`
	Address   Address  `json:"address"`
	Elevation float64  `json:"elevation"`
	Airly     bool     `json:"airly"`
	Sponsor   Sponsor  `json:"sponsor"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Address struct {
	Country         string `json:"country"`
	City            string `json:"city"`
	Street          string `json:"street"`
	Number          string `json:"number"`
	DisplayAddress1 string `json:"displayAddress1"`
	DisplayAddress2 string `json:"displayAddress2"`
}

type Sponsor struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Link        string `json:"link"`
	DisplayName string `json:"displayName"`
}

type Measurements struct {
	Current  Measurement   `json:"current"`
	History  []Measurement `json:"history"`
	Forecast []Measurement `json:"forecast"`
}

type Value struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type Index struct {
	Name        string  `json:"name"`
	Value       float64 `json:"value"`
	Level       string  `json:"level"`
	Description string  `json:"description"`
	Advice      string  `json:"advice"`
	Color       string  `json:"color"`
}

type Standard struct {
	Name      string  `json:"name"`
	Pollutant string  `json:"pollutant"`
	Limit     float64 `json:"limit"`
	Percent   float64 `json:"percent"`
}

type Measurement struct {
	FromDateTime time.Time  `json:"fromDateTime"`
	TillDateTime time.Time  `json:"tillDateTime"`
	Values       []Value    `json:"values"`
	Indexes      []Index    `json:"indexes"`
	Standards    []Standard `json:"standards"`
}

type IndexType struct {
	Name   string   `json:"name"`
	Levels []Levels `json:"levels"`
}
type Levels struct {
	Values      string `json:"values"`
	Level       string `json:"level"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type MeasurementType struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Unit  string `json:"unit"`
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	Key        string
	Language   string
	HttpClient HttpClient
}

const base = "https://airapi.airly.eu/v2/"

func (c Client) get(path string, v interface{}) error {
	req, err := http.NewRequest("GET", base+path, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("apikey", c.Key)
	if c.Language != "" {
		req.Header.Set("Accept-Language", c.Language)
	}
	client := c.HttpClient
	if client == nil {
		client = http.DefaultClient
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("%d: %s", res.StatusCode, body)
	}

	return json.Unmarshal(body, v)
}

func (c Client) Installation(id int) (Installation, error) {
	var i Installation
	err := c.get(fmt.Sprintf("installations/%d", id), &i)
	return i, err
}

func (c Client) NearestInstallations(loc Location, options ...NearestInstallationsOption) ([]Installation, error) {
	var i []Installation
	config := nearestInstallationsConfig{3.0, 1}
	for _, option := range options {
		option(&config)
	}
	err := c.get(fmt.Sprintf("installations/nearest?lat=%f&lng=%f&maxDistanceKM=%f&maxResults=%d",
		loc.Latitude, loc.Longitude, config.maxDistance, config.maxResults), &i)
	return i, err
}

func (c Client) NearestMeasurements(loc Location, options ...NearestInstallationsOption) (Measurements, error) {
	var m Measurements
	config := nearestInstallationsConfig{3.0, 1}
	for _, option := range options {
		option(&config)
	}
	err := c.get(fmt.Sprintf("measurements/nearest?lat=%f&lng=%f&maxDistanceKM=%f",
		loc.Latitude, loc.Longitude, config.maxDistance), &m)
	return m, err
}

func (c Client) PointMeasurements(loc Location) (Measurements, error) {
	var m Measurements
	err := c.get(fmt.Sprintf("measurements/point?lat=%f&lng=%f", loc.Latitude, loc.Longitude), &m)
	return m, err
}

func (c Client) InstallationMeasurements(installationId int) (Measurements, error) {
	var m Measurements
	err := c.get(fmt.Sprintf("measurements/installation?installationId=%d", installationId), &m)
	return m, err
}

type NearestInstallationsOption func(config *nearestInstallationsConfig)

func MaxDistance(maxDistance float64) NearestInstallationsOption {
	return func(c *nearestInstallationsConfig) {
		c.maxDistance = maxDistance
	}
}

func MaxResults(maxResults int) NearestInstallationsOption {
	return func(c *nearestInstallationsConfig) {
		c.maxResults = maxResults
	}
}

type nearestInstallationsConfig struct {
	maxDistance float64
	maxResults  int
}

func (c Client) IndexTypes() ([]IndexType, error) {
	var i []IndexType
	err := c.get("meta/measurements", &i)
	return i, err
}

func (c Client) MeasurementTypes() ([]MeasurementType, error) {
	var m []MeasurementType
	err := c.get("meta/measurements", &m)
	return m, err
}
