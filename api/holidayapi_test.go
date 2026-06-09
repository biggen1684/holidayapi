package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Testing Listcountries function to see if valid data is marshaled correctly
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

// Testing that wrong endpoint is returned as 404 error
func TestListCountries_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()
	client := server.Client()
	countries, err := ListCountries(client, server.URL, false)
	assert.Nil(t, countries)
	assert.EqualError(t, err, "countries endpoint not found at URL")
}

// Test GetHolidays if valid data is marshaled correctly
func TestGetHolidays(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"date":"2026-01-01","localName":"New Year's Day","name":"New Year's Day","countryCode":"US","fixed":false,"global":true,"counties":null,"launchYear":null,"types":["Public","Bank"]}]`))
	}))
	defer server.Close()
	client := server.Client()
	holidays, err := GetHolidays(client, server.URL, "2026", "US", false)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(holidays))
	assert.Equal(t, "New Year's Day", holidays[0].Name)
}

// Test GetHolidays if an unsupported is handled correctly
func TestGetHolidays_Bad_Year(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()
	client := server.Client()
	holidays, err := GetHolidays(client, server.URL, "1000", "US", false)
	assert.Nil(t, holidays)
	assert.Error(t, err)
}

func TestGetHolidays_Bad_Country(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()
	client := server.Client()
	holidays, err := GetHolidays(client, server.URL, "2000", "ZZZ", false)
	assert.Nil(t, holidays)
	assert.Error(t, err)
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
