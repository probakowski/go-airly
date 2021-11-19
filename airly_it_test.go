package airly

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInstallationIT(t *testing.T) {
	apiKey := os.Getenv("apikey")
	if apiKey == "" {
		t.Skip("skipping testing - no apikey provided")
		return
	}
	if testing.Short() {
		t.Skip("skipping testing in short mode")
		return
	}
	api := Client{
		Key:      apiKey,
		Language: "pl",
	}
	installation, err := api.Installation(8077)
	assert.Nil(t, err)
	assert.Equal(t, Installation{
		Id:        8077,
		Location:  Location{50.062006, 19.940984},
		Address:   Address{"Poland", "Krakow", "Mikołajska", "", "Krakow", "Mikołajska"},
		Elevation: 220.38,
		Airly:     true,
		Sponsor:   Sponsor{489, "Chatham Financial", "sponsor sensora Airly", "https://cdn.airly.eu/logo/ChathamFinancial_1570109001008_473803190.jpg", "https://crossweb.pl/job/chatham-financial/ ", "Chatham Financial"},
	}, installation)
}
