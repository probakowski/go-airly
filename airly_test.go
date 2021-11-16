package airly

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type mockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func readCloser(s string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(s))
}

func (m mockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestError(t *testing.T) {
	err := errors.New("error")
	api := Api{
		Key:      "x1234x",
		Language: "pl",
		HttpClient: mockClient{func(req *http.Request) (*http.Response, error) {
			return nil, err
		}}}
	_, err2 := api.Installation(204)
	assert.Equal(t, err, err2)
}

func TestNon200Status(t *testing.T) {
	api := Api{
		Key:      "x1234x",
		Language: "pl",
		HttpClient: mockClient{func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 404,
				Body:       readCloser("not found"),
			}, nil
		}}}
	_, err2 := api.Installation(204)
	assert.Equal(t, "404: not found", err2.Error())
}

func TestInstallation(t *testing.T) {
	api := Api{
		Key:      "x1234x",
		Language: "pl",
		HttpClient: mockClient{func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "https://airapi.airly.eu/v2/installations/204", req.URL.String())
			assert.Equal(t, "application/json", req.Header.Get("Accept"))
			assert.Equal(t, "pl", req.Header.Get("Accept-Language"))
			assert.Equal(t, "x1234x", req.Header.Get("apikey"))
			return &http.Response{
				StatusCode: 200,
				Body: readCloser(`{
									  "id": 204,
									  "location": {
										"latitude": 50.062006,
										"longitude": 19.940984
									  },
									  "address": {
										"country": "Poland",
										"city": "Kraków",
										"street": "Mikołajska",
										"number": "4B",
										"displayAddress1": "Kraków",
										"displayAddress2": "Mikołajska"
									  },
									  "elevation": 220.38,
									  "airly": true,
									  "sponsor": {
										"name": "KrakówOddycha",
										"description": "Sensor Airly w ramach akcji",
										"logo": "https://cdn.airly.org/logo/KrakówOddycha.jpg",
										"link": "https://przykladowy_link_do_strony_sponsora.pl"
									  }
									}`),
			}, nil
		}},
	}
	installation, err := api.Installation(204)
	assert.Nil(t, err)
	assert.Equal(t, Installation{
		Id: 204,
		Location: Location{
			Latitude:  50.062006,
			Longitude: 19.940984,
		},
		Address: Address{
			Country:         "Poland",
			City:            "Kraków",
			Street:          "Mikołajska",
			Number:          "4B",
			DisplayAddress1: "Kraków",
			DisplayAddress2: "Mikołajska",
		},
		Elevation: 220.38,
		Airly:     true,
		Sponsor: Sponsor{
			Name:        "KrakówOddycha",
			Description: "Sensor Airly w ramach akcji",
			Logo:        "https://cdn.airly.org/logo/KrakówOddycha.jpg",
			Link:        "https://przykladowy_link_do_strony_sponsora.pl",
		},
	}, installation)
}

func TestNearestInstallations(t *testing.T) {
	api := Api{
		Key:      "x1234x",
		Language: "pl",
		HttpClient: mockClient{func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "https://airapi.airly.eu/v2/installations/nearest?lat=50.062006"+
				"&lng=19.940984&maxDistanceKM=3.000000&maxResults=1", req.URL.String())
			assert.Equal(t, "application/json", req.Header.Get("Accept"))
			assert.Equal(t, "pl", req.Header.Get("Accept-Language"))
			assert.Equal(t, "x1234x", req.Header.Get("apikey"))
			return &http.Response{
				StatusCode: 200,
				Body: readCloser(`[{
									  "id": 204,
									  "location": {
										"latitude": 50.062006,
										"longitude": 19.940984
									  },
									  "address": {
										"country": "Poland",
										"city": "Kraków",
										"street": "Mikołajska",
										"number": "4B",
										"displayAddress1": "Kraków",
										"displayAddress2": "Mikołajska"
									  },
									  "elevation": 220.38,
									  "airly": true,
									  "sponsor": {
										"name": "KrakówOddycha",
										"description": "Sensor Airly w ramach akcji",
										"logo": "https://cdn.airly.org/logo/KrakówOddycha.jpg",
										"link": "https://przykladowy_link_do_strony_sponsora.pl"
									  }
									}]`),
			}, nil
		}},
	}
	installations, err := api.NearestInstallations(Location{50.062006, 19.940984})
	assert.Nil(t, err)
	assert.Equal(t, []Installation{{
		Id: 204,
		Location: Location{
			Latitude:  50.062006,
			Longitude: 19.940984,
		},
		Address: Address{
			Country:         "Poland",
			City:            "Kraków",
			Street:          "Mikołajska",
			Number:          "4B",
			DisplayAddress1: "Kraków",
			DisplayAddress2: "Mikołajska",
		},
		Elevation: 220.38,
		Airly:     true,
		Sponsor: Sponsor{
			Name:        "KrakówOddycha",
			Description: "Sensor Airly w ramach akcji",
			Logo:        "https://cdn.airly.org/logo/KrakówOddycha.jpg",
			Link:        "https://przykladowy_link_do_strony_sponsora.pl",
		},
	}}, installations)
}

func TestNearestInstallationsOptions(t *testing.T) {
	api := Api{
		Key:      "x1234x",
		Language: "pl",
		HttpClient: mockClient{func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "https://airapi.airly.eu/v2/installations/nearest?lat=50.062006"+
				"&lng=19.940984&maxDistanceKM=5.000000&maxResults=3", req.URL.String())
			assert.Equal(t, "application/json", req.Header.Get("Accept"))
			assert.Equal(t, "pl", req.Header.Get("Accept-Language"))
			assert.Equal(t, "x1234x", req.Header.Get("apikey"))
			return &http.Response{
				StatusCode: 200,
				Body: readCloser(`[{
									  "id": 204,
									  "location": {
										"latitude": 50.062006,
										"longitude": 19.940984
									  },
									  "address": {
										"country": "Poland",
										"city": "Kraków",
										"street": "Mikołajska",
										"number": "4B",
										"displayAddress1": "Kraków",
										"displayAddress2": "Mikołajska"
									  },
									  "elevation": 220.38,
									  "airly": true,
									  "sponsor": {
										"name": "KrakówOddycha",
										"description": "Sensor Airly w ramach akcji",
										"logo": "https://cdn.airly.org/logo/KrakówOddycha.jpg",
										"link": "https://przykladowy_link_do_strony_sponsora.pl"
									  }
									}]`),
			}, nil
		}},
	}
	installations, err := api.NearestInstallations(Location{50.062006, 19.940984},
		MaxDistance(5), MaxResults(3))
	assert.Nil(t, err)
	assert.Equal(t, []Installation{{
		Id: 204,
		Location: Location{
			Latitude:  50.062006,
			Longitude: 19.940984,
		},
		Address: Address{
			Country:         "Poland",
			City:            "Kraków",
			Street:          "Mikołajska",
			Number:          "4B",
			DisplayAddress1: "Kraków",
			DisplayAddress2: "Mikołajska",
		},
		Elevation: 220.38,
		Airly:     true,
		Sponsor: Sponsor{
			Name:        "KrakówOddycha",
			Description: "Sensor Airly w ramach akcji",
			Logo:        "https://cdn.airly.org/logo/KrakówOddycha.jpg",
			Link:        "https://przykladowy_link_do_strony_sponsora.pl",
		},
	}}, installations)
}

func TestInstallationMeasurements(t *testing.T) {
	api := Api{
		Key:      "x1234x",
		Language: "pl",
		HttpClient: mockClient{func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "https://airapi.airly.eu/v2/measurements/installation?installationId=204", req.URL.String())
			assert.Equal(t, "application/json", req.Header.Get("Accept"))
			assert.Equal(t, "pl", req.Header.Get("Accept-Language"))
			assert.Equal(t, "x1234x", req.Header.Get("apikey"))
			return &http.Response{
				StatusCode: 200,
				Body: readCloser(`{
									  "current": {
										"fromDateTime": "2018-08-24T08:24:48.652Z",
										"tillDateTime": "2018-08-24T09:24:48.652Z",
										"values": [
										  { "name": "PM1",          "value": 12.73   },
										  { "name": "PM25",         "value": 18.7    }
										],
										"indexes": [
										  {
											"name": "AIRLY_CAQI",
											"value": 35.53,
											"level": "LOW",
											"description": "Dobre powietrze.",
											"advice": "Możesz bez obaw wyjść na zewnątrz.",
											"color": "#D1CF1E"
										  }
										],
										"standards": [
										  {
											"name": "WHO",
											"pollutant": "PM25",
											"limit": 25,
											"percent": 74.81
										  }
										]
									  },
									  "history": [],
									  "forecast": []
									}`),
			}, nil
		}},
	}
	measurements, err := api.InstallationMeasurements(204)
	assert.Nil(t, err)
	assert.Equal(t, Measurements{
		Current: Measurement{
			FromDateTime: time.Date(2018, 8, 24, 8, 24, 48, 652*1000*1000, time.UTC),
			TillDateTime: time.Date(2018, 8, 24, 9, 24, 48, 652*1000*1000, time.UTC),
			Values: []Value{{
				Name:  "PM1",
				Value: 12.73,
			}, {
				Name:  "PM25",
				Value: 18.7,
			}},
			Indexes: []Index{{
				Name:        "AIRLY_CAQI",
				Value:       35.53,
				Level:       "LOW",
				Description: "Dobre powietrze.",
				Advice:      "Możesz bez obaw wyjść na zewnątrz.",
				Color:       "#D1CF1E",
			}},
			Standards: []Standard{{
				Name:      "WHO",
				Pollutant: "PM25",
				Limit:     25,
				Percent:   74.81,
			}},
		},
		History:  []Measurement{},
		Forecast: []Measurement{},
	}, measurements)
}

func TestNearestMeasurements(t *testing.T) {
	api := Api{
		Key:      "x1234x",
		Language: "pl",
		HttpClient: mockClient{func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "https://airapi.airly.eu/v2/measurements/nearest?lat=50.062006&lng=19.940984&maxDistanceKM=5.000000", req.URL.String())
			assert.Equal(t, "application/json", req.Header.Get("Accept"))
			assert.Equal(t, "pl", req.Header.Get("Accept-Language"))
			assert.Equal(t, "x1234x", req.Header.Get("apikey"))
			return &http.Response{
				StatusCode: 200,
				Body: readCloser(`{
									  "current": {
										"fromDateTime": "2018-08-24T08:24:48.652Z",
										"tillDateTime": "2018-08-24T09:24:48.652Z",
										"values": [
										  { "name": "PM1",          "value": 12.73   },
										  { "name": "PM25",         "value": 18.7    }
										],
										"indexes": [
										  {
											"name": "AIRLY_CAQI",
											"value": 35.53,
											"level": "LOW",
											"description": "Dobre powietrze.",
											"advice": "Możesz bez obaw wyjść na zewnątrz.",
											"color": "#D1CF1E"
										  }
										],
										"standards": [
										  {
											"name": "WHO",
											"pollutant": "PM25",
											"limit": 25,
											"percent": 74.81
										  }
										]
									  },
									  "history": [],
									  "forecast": []
									}`),
			}, nil
		}},
	}
	measurements, err := api.NearestMeasurements(Location{50.062006, 19.940984}, MaxDistance(5))
	assert.Nil(t, err)
	assert.Equal(t, Measurements{
		Current: Measurement{
			FromDateTime: time.Date(2018, 8, 24, 8, 24, 48, 652*1000*1000, time.UTC),
			TillDateTime: time.Date(2018, 8, 24, 9, 24, 48, 652*1000*1000, time.UTC),
			Values: []Value{{
				Name:  "PM1",
				Value: 12.73,
			}, {
				Name:  "PM25",
				Value: 18.7,
			}},
			Indexes: []Index{{
				Name:        "AIRLY_CAQI",
				Value:       35.53,
				Level:       "LOW",
				Description: "Dobre powietrze.",
				Advice:      "Możesz bez obaw wyjść na zewnątrz.",
				Color:       "#D1CF1E",
			}},
			Standards: []Standard{{
				Name:      "WHO",
				Pollutant: "PM25",
				Limit:     25,
				Percent:   74.81,
			}},
		},
		History:  []Measurement{},
		Forecast: []Measurement{},
	}, measurements)
}

func TestPointMeasurements(t *testing.T) {
	api := Api{
		Key:      "x1234x",
		Language: "pl",
		HttpClient: mockClient{func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "https://airapi.airly.eu/v2/measurements/point?lat=50.062006&lng=19.940984", req.URL.String())
			assert.Equal(t, "application/json", req.Header.Get("Accept"))
			assert.Equal(t, "pl", req.Header.Get("Accept-Language"))
			assert.Equal(t, "x1234x", req.Header.Get("apikey"))
			return &http.Response{
				StatusCode: 200,
				Body: readCloser(`{
									  "current": {
										"fromDateTime": "2018-08-24T08:24:48.652Z",
										"tillDateTime": "2018-08-24T09:24:48.652Z",
										"values": [
										  { "name": "PM1",          "value": 12.73   },
										  { "name": "PM25",         "value": 18.7    }
										],
										"indexes": [
										  {
											"name": "AIRLY_CAQI",
											"value": 35.53,
											"level": "LOW",
											"description": "Dobre powietrze.",
											"advice": "Możesz bez obaw wyjść na zewnątrz.",
											"color": "#D1CF1E"
										  }
										],
										"standards": [
										  {
											"name": "WHO",
											"pollutant": "PM25",
											"limit": 25,
											"percent": 74.81
										  }
										]
									  },
									  "history": [],
									  "forecast": []
									}`),
			}, nil
		}},
	}
	measurements, err := api.PointMeasurements(Location{50.062006, 19.940984})
	assert.Nil(t, err)
	assert.Equal(t, Measurements{
		Current: Measurement{
			FromDateTime: time.Date(2018, 8, 24, 8, 24, 48, 652*1000*1000, time.UTC),
			TillDateTime: time.Date(2018, 8, 24, 9, 24, 48, 652*1000*1000, time.UTC),
			Values: []Value{{
				Name:  "PM1",
				Value: 12.73,
			}, {
				Name:  "PM25",
				Value: 18.7,
			}},
			Indexes: []Index{{
				Name:        "AIRLY_CAQI",
				Value:       35.53,
				Level:       "LOW",
				Description: "Dobre powietrze.",
				Advice:      "Możesz bez obaw wyjść na zewnątrz.",
				Color:       "#D1CF1E",
			}},
			Standards: []Standard{{
				Name:      "WHO",
				Pollutant: "PM25",
				Limit:     25,
				Percent:   74.81,
			}},
		},
		History:  []Measurement{},
		Forecast: []Measurement{},
	}, measurements)
}
