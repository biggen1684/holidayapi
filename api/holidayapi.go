package api

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Holiday struct {
	Date        string `json:"date"`
	CountryCode string `json:"countryCode"`
	Name        string `json:"name"`
	Global      bool   `json:"global"`
}

func GetHolidays() ([]Holiday, error) {
	currentYear := fmt.Sprintf("%d", time.Now().Year())
	year := flag.String("year", currentYear, "the year in xxxx format")
	countryCode := flag.String("countrycode", "US", "2-letter ISO 3166-1 alpha-2 country code")
	debug := flag.Bool("debug", false, "print raw API response")
	flag.Parse()
	url := fmt.Sprintf("https://date.nager.at/api/v3/PublicHolidays/%s/%s", *year, *countryCode)
	client := &http.Client{Timeout: 30 * time.Second}

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
		return nil, fmt.Errorf("invalid country code: %q", *countryCode)
	}
	if res.StatusCode == http.StatusBadRequest {
		return nil, fmt.Errorf("invalid year: %q", *year)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %s", err)
	}

	if *debug {
		fmt.Println(string(body))
	}

	var holidays []Holiday
	err = json.Unmarshal(body, &holidays)
	if err != nil {
		return nil, fmt.Errorf("decoding body: %s", err)
	}

	return holidays, nil
}
