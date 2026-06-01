package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	url := "https://date.nager.at/api/v3/AvailableCountries"
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("connecting to API: %s", err)
	}

	req.Header.Add("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("getting request: %s", err)
	}
	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("countries endpoint not found at URL")
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %s", err)
	}

	if debug == true {
		fmt.Println(string(body))
	}

	var countries []Countries
	err = json.Unmarshal(body, &countries)
	if err != nil {
		return nil, fmt.Errorf("decoding body: %s", err)
	}
	return countries, nil

}

func GetHolidays(client *http.Client, year string, countryCode string, debug bool) ([]Holiday, error) {

	url := fmt.Sprintf("https://date.nager.at/api/v3/PublicHolidays/%s/%s", year, countryCode)
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("connecting to API: %s", err)
	}

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

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %s", err)
	}

	if debug == true {
		fmt.Println(string(body))
	}

	var holidays []Holiday
	err = json.Unmarshal(body, &holidays)
	if err != nil {
		return nil, fmt.Errorf("decoding body: %s", err)
	}

	return holidays, nil
}
