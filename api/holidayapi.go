package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Holiday struct {
	Date        string `json:"date"`
	CountryCode string `json:"countryCode"`
	Name        string `json:"name"`
	Global      bool   `json:"global"`
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

	//Request Header and deal with errors if needed
	req.Header.Add("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("getting request: %s", err)
	}
	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("countries endpoint not found at URL")
	}

	//Read the raw body and close when done
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %s", err)
	}

	//Print raw body if debug flag true
	if debug == true {
		fmt.Println(string(body))
	}

	//Finally unmarshal into a slice containing the struct defined
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
	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("invalid country code: %q", countryCode)
	}
	if res.StatusCode == http.StatusBadRequest {
		return nil, fmt.Errorf("invalid year: %q", year)
	}

	//Read the raw body and close when done
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %s", err)
	}

	//Print raw body if debug flag true
	if debug == true {
		fmt.Println(string(body))
	}

	//Finally unmarshal into a slice containing the struct defined
	var holidays []Holiday
	err = json.Unmarshal(body, &holidays)
	if err != nil {
		return nil, fmt.Errorf("decoding body: %s", err)
	}

	return holidays, nil
}

func PrintCountries(countries []Countries) {
	fmt.Print("The two letter codes for all countries are as follows.\n\n")
	for i, v := range countries {
		fmt.Printf("%d. Two letter code for %s is '%s'.\n", i+1, v.Name, v.Code)
	}
}

func PrintHolidays(holidays []Holiday, year string, countryCode string, globalOnly bool) {
	fmt.Printf("The holidays in %s for country %s are as follows:\n\n", year, countryCode)

	if globalOnly == true {
		count := 1
		for _, v := range holidays {
			if v.Global == true {
				fmt.Printf("%d. %s %s\n", count, strings.TrimPrefix(v.Date, year+"-"), v.Name)
				count++
			}
		}
		return
	}
	for i, v := range holidays {
		fmt.Printf("%d. %s %s\n", i+1, strings.TrimPrefix(v.Date, year+"-"), v.Name)
	}
}
