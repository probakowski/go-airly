// Package airly provides wrapper for Airly API https://developer.airly.org/docs
package airly

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Installation metadata, see https://developer.airly.org/docs#endpoints.installations
type Installation struct {
	Id        int      `json:"id"`
	Location  Location `json:"location"`
	Address   Address  `json:"address"`
	Elevation float64  `json:"elevation"`
	Airly     bool     `json:"airly"`
	Sponsor   Sponsor  `json:"sponsor"`
}

// Location represents geographical location given by coordinates, used for Nearest* APIs
// and as part of Installation metadata
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Address of Installation
type Address struct {
	Country         string `json:"country"`
	City            string `json:"city"`
	Street          string `json:"street"`
	Number          string `json:"number"`
	DisplayAddress1 string `json:"displayAddress1"`
	DisplayAddress2 string `json:"displayAddress2"`
}

// Sponsor of Installation, see https://developer.airly.org/docs#concepts.installations.sponsors
type Sponsor struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Link        string `json:"link"`
	DisplayName string `json:"displayName"`
}

// Measurements data, see https://developer.airly.org/docs#endpoints.measurements
type Measurements struct {
	Current  Measurement   `json:"current"`
	History  []Measurement `json:"history"`
	Forecast []Measurement `json:"forecast"`
}

// Value of measurement
type Value struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// Index showing aggregated air quality
type Index struct {
	Name        string  `json:"name"`
	Value       float64 `json:"value"`
	Level       string  `json:"level"`
	Description string  `json:"description"`
	Advice      string  `json:"advice"`
	Color       string  `json:"color"`
}

// Standard used for measuring
type Standard struct {
	Name      string  `json:"name"`
	Pollutant string  `json:"pollutant"`
	Limit     float64 `json:"limit"`
	Percent   float64 `json:"percent"`
}

// Measurement represents aggregated values
type Measurement struct {
	FromDateTime time.Time  `json:"fromDateTime"`
	TillDateTime time.Time  `json:"tillDateTime"`
	Values       []Value    `json:"values"`
	Indexes      []Index    `json:"indexes"`
	Standards    []Standard `json:"standards"`
}

// IndexType represents index metadata, https://developer.airly.org/docs#endpoints.meta.indexes
type IndexType struct {
	Name   string  `json:"name"`
	Levels []Level `json:"levels"`
}

// Level metadata
type Level struct {
	Values      string `json:"values"`
	Level       string `json:"level"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

// MeasurementType metadata, see https://developer.airly.org/docs#endpoints.meta.measurements
type MeasurementType struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Unit  string `json:"unit"`
}

// HttpClient to use for requests
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client for Airly API
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

// Installation returns installation by id. See https://developer.airly.org/docs#endpoints.installations.getbyid
func (c Client) Installation(id int) (Installation, error) {
	var i Installation
	err := c.get(fmt.Sprintf("installations/%d", id), &i)
	return i, err
}

// NearestInstallations returns installations near specified point, range can be defined with MaxDistance,
// number of results can be defined with MaxResults. See https://developer.airly.org/docs#endpoints.installations.nearest
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

// NearestMeasurements returns measurements for an installation closest to a given location, range can be defined with MaxDistance.
// See https://developer.airly.org/en/docs#endpoints.measurements.nearest
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

// PointMeasurements returns any geographical location.
// Measurement values are interpolated by averaging measurements from nearby sensors (up to 1,5km away from the given point).
// The returned value is a weighted average, with the weight inversely proportional to the distance from the sensor to the given point.
// See https://developer.airly.org/docs#endpoints.measurements.point
func (c Client) PointMeasurements(loc Location) (Measurements, error) {
	var m Measurements
	err := c.get(fmt.Sprintf("measurements/point?lat=%f&lng=%f", loc.Latitude, loc.Longitude), &m)
	return m, err
}

// InstallationMeasurements returns measurements for concrete installation, see https://developer.airly.org/docs#endpoints.measurements.installation
func (c Client) InstallationMeasurements(installationId int) (Measurements, error) {
	var m Measurements
	err := c.get(fmt.Sprintf("measurements/installation?installationId=%d", installationId), &m)
	return m, err
}

// NearestInstallationsOption represents option to narrow search results
type NearestInstallationsOption func(config *nearestInstallationsConfig)

// MaxDistance to given points in km
func MaxDistance(maxDistance float64) NearestInstallationsOption {
	return func(c *nearestInstallationsConfig) {
		c.maxDistance = maxDistance
	}
}

// MaxResults that can be returned by API call
func MaxResults(maxResults int) NearestInstallationsOption {
	return func(c *nearestInstallationsConfig) {
		c.maxResults = maxResults
	}
}

type nearestInstallationsConfig struct {
	maxDistance float64
	maxResults  int
}

// IndexTypes returns a list of all the index types supported in the API along with lists of levels defined
// per each index type, see https://developer.airly.org/docs#endpoints.meta.indexes
func (c Client) IndexTypes() ([]IndexType, error) {
	var i []IndexType
	err := c.get("meta/measurements", &i)
	return i, err
}

// MeasurementTypes returns list of all the measurement types supported in the API along with their names and units,
// see https://developer.airly.org/docs#endpoints.meta.measurements
func (c Client) MeasurementTypes() ([]MeasurementType, error) {
	var m []MeasurementType
	err := c.get("meta/measurements", &m)
	return m, err
}
