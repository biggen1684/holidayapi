package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	//Parse date returned from API to actual day of week for printing later
	for i, v := range holidays {
		t, err := time.Parse("2006-01-02", v.Date)
		if err != nil {
			return nil, fmt.Errorf("problem parsing date to weekday %s", err)
		}
		holidays[i].Weekday = t.Weekday().String()
	}

	return holidays, nil
}

// Prints all available countries if flag is used
func PrintCountries(countries []Countries) {
	fmt.Print("The two letter codes for all countries are as follows.\n\n")
	for i, v := range countries {
		fmt.Printf("%d. Two letter code for %s is '%s'.\n", i+1, v.Name, v.Code)
	}
}

// Prints holidays with a simple counter and the year removed from print line
// Default to printing only global holidays (national holidays)
// If "-globalonly=false" flag is passed in, we print all known holidays
func PrintHolidays(holidays []Holiday, year string, countryCode string, globalOnly bool) {
	fmt.Printf("The holidays in %s for the country of %s are as follows:\n\n", year, countryCode)
	count := 1
	for _, v := range holidays {
		if globalOnly && !v.Global {
			continue
		}
		fmt.Printf("%d. %s, %s %s\n", count, v.Weekday, strings.TrimPrefix(v.Date, year+"-"), v.Name)
		count++
	}
}
