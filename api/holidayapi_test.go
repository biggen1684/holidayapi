package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Testing Listcountries function to see if data is marshaled correctly
func TestListCountries(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"countryCode":"US","name":"United States"},{"countryCode":"CA","name":"Canada"}]`))
	}))
	defer server.Close()
	client := server.Client()
	countries, err := ListCountries(client, server.URL, false)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(countries))
	assert.Equal(t, "United States", countries[0].Name)
	assert.Equal(t, "US", countries[0].Code)
	assert.Equal(t, "Canada", countries[1].Name)
	assert.Equal(t, "CA", countries[1].Code)
}

func TestWeekday(t *testing.T) {
	date := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, "Thursday", weekday(date))
}

func TestUnderThirty_False(t *testing.T) {
	date := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	assert.False(t, underThirty(date))
}

func TestUnderThirty_True(t *testing.T) {
	date := time.Now().Add(24 * time.Hour)
	assert.True(t, underThirty(date))
}

// Truncate to midnight so we calculate whole days from start of today
// rather than from the exact current time, which would cause off-by-one errors
func TestDaysAway(t *testing.T) {
	n := time.Now()
	date := time.Date(n.Year(), n.Month(), n.Day()+5, 0, 0, 0, 0, time.Local)
	assert.Equal(t, 5, daysAway(date))
}

func TestIsToday(t *testing.T) {
	today := time.Now()
	date := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.Local)
	assert.True(t, isToday(date))
}
