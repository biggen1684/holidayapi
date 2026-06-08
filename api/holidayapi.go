package api

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

type Holiday struct {
	Date        string `json:"date"`
	CountryCode string `json:"countryCode"`
	Name        string `json:"name"`
	Global      bool   `json:"global"`
	Weekday     string `json:"weekday"`
	UnderThirty bool   `json:"underThirty"`
	DaysAway    int    `json:"daysAway"`
	IsToday     bool   `json:"isToday"`
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
	// Parse date returned from API to actual day of week for printing
	// Use same loop for other enrichment functions
	for i, v := range rawHolidays {
		t, err := time.ParseInLocation("2006-01-02", v.Date, time.Local)
		if err != nil {
			return nil, fmt.Errorf("problem parsing date to time.Time %s", err)
		}
		rawHolidays[i].Weekday = weekday(t)
		rawHolidays[i].UnderThirty = underThirty(t)
		rawHolidays[i].DaysAway = daysAway(t)
		rawHolidays[i].IsToday = isToday(t)

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
// Negative numbers indicate holiday is in the past
func daysAway(t time.Time) int {
	n := time.Now()
	now := time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.Local)
	timeUntil := int(math.Round(t.Sub(now).Hours() / 24))
	return timeUntil
}

// Determine if today is a holiday and return true if so
func isToday(t time.Time) bool {
	now := time.Now()
	if now.Year() == t.Year() && now.Month() == t.Month() && now.Day() == t.Day() {
		return true
	}
	return false

}

// Output holidays with enrichments to .csv file with a header
func OutputCSV(holidays []Holiday) error {
	file, err := os.Create("holidays.csv")
	if err != nil {
		return fmt.Errorf("failed creating file: %s", err)
	}
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()
	header := []string{"Date", "CountryCode", "Name", "Global", "Weekday", "UnderThirty", "DaysAway", "isToday"}
	err = w.Write(header)
	if err != nil {
		return fmt.Errorf("problem writing header: %s", err)
	}
	for _, holiday := range holidays {
		row := []string{holiday.Date, holiday.CountryCode, holiday.Name,
			strconv.FormatBool(holiday.Global), holiday.Weekday,
			strconv.FormatBool(holiday.UnderThirty),
			strconv.Itoa(holiday.DaysAway),
			strconv.FormatBool(holiday.IsToday),
		}
		err = w.Write(row)
		if err != nil {
			return fmt.Errorf("problem writing to file: %s", err)
		}

	}
	fmt.Println("holiday.csv file saved successfully")
	return nil
}

// Output holidays with enrichments to .json file
func OutputJSON(holidays []Holiday) error {
	file, err := os.Create("holidays.json")
	if err != nil {
		return fmt.Errorf("failed creating file: %s", err)
	}
	defer file.Close()
	jsonData, err := json.Marshal(holidays)
	if err != nil {
		return fmt.Errorf("problem running marshaler: %s", err)
	}
	err = os.WriteFile("holidays.json", jsonData, 0644)
	if err != nil {
		return fmt.Errorf("problem writing file: %s", err)
	}
	fmt.Println("holiday.json file saved successfully")
	return nil
}

// Prints all available countries if flag is used using text/tabwriter
func PrintCountries(countries []Countries) {
	fmt.Print("The two letter codes for all countries are as follows.\n\n")
	// for i, v := range countries {
	// 	fmt.Printf("%d. Two letter code for %s is '%s'.\n", i+1, v.Name, v.Code)
	// }
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	for i, v := range countries {
		fmt.Fprintf(w, "%d. %s\t", i+1, v.Name)
		if (i+1)%5 == 0 {
			fmt.Fprintln(w)
		}
	}
	w.Flush()
	fmt.Println()
}

// Prints holidays with a simple counter and the year removed from print line
// Default to printing only federal holidays (national holidays)
// If "-federalonly=false" flag is passed in, we print all known holidays
// Also prints holidays as blue if under 30 days away and color=true (default) flag
// Print out a statment if today is a holiday at the top
func PrintHolidays(holidays []Holiday, year string, countryCode string, federalOnly bool, color bool) {
	fmt.Printf("The holidays in %s for the country of %s are as follows:\n\n", year, countryCode)
	for _, v := range holidays {
		if v.IsToday {
			fmt.Printf("🎉 Today is %s!\n\n", v.Name)
		}
	}
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
