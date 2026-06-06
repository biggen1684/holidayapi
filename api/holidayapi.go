package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"
)

type Holiday struct {
	Date        string `json:"date"`
	CountryCode string `json:"countryCode"`
	Name        string `json:"name"`
	Global      bool   `json:"global"`
	Weekday     string
	UnderThirty bool
	DaysAway    int
}

type Countries struct {
	Code string `json:"countryCode"`
	Name string `json:"name"`
}

func ListCountries(client *http.Client, debug bool) ([]Countries, error) {

	//Setup context, Get, and URL
	url := "https://date.nager.at/api/v3/AvailableCountries"
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("connecting to API: %s", err)
	}

	//Request Header and deal with errors if needed and close when function ends
	req.Header.Add("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("getting request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("countries endpoint not found at URL")
	}

	//Read the raw body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %s", err)
	}

	//Print raw body if debug flag true
	if debug == true {
		fmt.Println(string(body))
	}

	//Finally unmarshal into a slice containing the struct declared above
	var countries []Countries
	err = json.Unmarshal(body, &countries)
	if err != nil {
		return nil, fmt.Errorf("decoding body: %s", err)
	}
	return countries, nil

}

func GetHolidays(client *http.Client, year string, countryCode string, debug bool) ([]Holiday, error) {

	//Setup context, Get, and URL
	url := fmt.Sprintf("https://date.nager.at/api/v3/PublicHolidays/%s/%s", year, countryCode)
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("connecting to API: %s", err)
	}

	//Request Header and deal with errors if needed
	req.Header.Add("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("getting request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("invalid country code: %q", countryCode)
	}
	if res.StatusCode == http.StatusBadRequest {
		return nil, fmt.Errorf("invalid year: %q", year)
	}

	//Read the raw body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %s", err)
	}

	//Print raw body if debug flag true
	if debug == true {
		fmt.Println(string(body))
	}

	//Finally unmarshal into a slice containing the struct declared above
	var holidays []Holiday
	err = json.Unmarshal(body, &holidays)
	if err != nil {
		return nil, fmt.Errorf("decoding body: %s", err)
	}
	return holidays, nil
}

func EnrichHolidays(rawHolidays []Holiday) ([]Holiday, error) {
	//Parse date returned from API to actual day of week for printing
	//Use same loop to set the Under30 days bool to true for color printing
	for i, v := range rawHolidays {
		t, err := time.ParseInLocation("2006-01-02", v.Date, time.Local)
		if err != nil {
			return nil, fmt.Errorf("problem parsing date to time.Time %s", err)
		}
		rawHolidays[i].Weekday = weekday(t)
		rawHolidays[i].UnderThirty = underThirty(t)
		rawHolidays[i].DaysAway = daysaway(t)

	}
	return rawHolidays, nil
}

// Helper function that returns weekday in string format
func weekday(t time.Time) string {
	return t.Weekday().String()
}

// Helper function that returns true if holiday is < 30 days away
func underThirty(t time.Time) bool {
	if time.Until(t) < 30*24*time.Hour && time.Until(t) > 0 {
		return true
	}
	return false
}

// Helper function that calculates how far holiday is (in days) away
// Set today's date to midnight in local time zone
// Then calculate difference from today and t rounding to nearest whole day
// Retun 0 if negative (date has passed)
func daysaway(t time.Time) int {
	n := time.Now()
	now := time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.Local)
	timeUntil := int(math.Round(t.Sub(now).Hours() / 24))
	if timeUntil > 0 {
		return timeUntil
	}
	return 0
}

// Prints all available countries if flag is used
func PrintCountries(countries []Countries) {
	fmt.Print("The two letter codes for all countries are as follows.\n\n")
	for i, v := range countries {
		fmt.Printf("%d. Two letter code for %s is '%s'.\n", i+1, v.Name, v.Code)
	}
}

// Prints holidays with a simple counter and the year removed from print line
// Default to printing only federal holidays (national holidays)
// If "-federalonly=false" flag is passed in, we print all known holidays
// Also prints holidays as blue if under 30 days away and color=true (default) flag
func PrintHolidays(holidays []Holiday, year string, countryCode string, federalOnly bool, color bool) {
	fmt.Printf("The holidays in %s for the country of %s are as follows:\n\n", year, countryCode)
	count := 1
	for _, v := range holidays {
		if federalOnly && !v.Global {
			continue
		}
		daysAway := ""
		if v.DaysAway > 0 {
			daysAway = fmt.Sprintf("(%d days away)", v.DaysAway)
		}
		line := fmt.Sprintf("%d. %s, %s %s %s\n", count, v.Weekday, strings.TrimPrefix(v.Date, year+"-"), v.Name, daysAway)
		if color && v.UnderThirty {
			fmt.Printf("\033[1;34m%s\033[0m", line) // blue ANSI = under 30 days
		} else {
			fmt.Print(line)
		}
		count++
	}
}
