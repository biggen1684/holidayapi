package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/biggen1684/holidayapi/api"
)

func main() {

	client := &http.Client{Timeout: 30 * time.Second}

	// Create flags to pass into CLI
	currentYear := fmt.Sprintf("%d", time.Now().Year())
	year := flag.String("year", currentYear, "the year in YYYY format")
	countryCode := flag.String("countrycode", "US", "2-letter ISO 3166-1 alpha-2 country code")
	debug := flag.Bool("debug", false, "print raw API response, use -debug to enable")
	listCountries := flag.Bool("listcountries", false, "list all available countries and their 2 letter codes, use -listcountries to enable")
	federalOnly := flag.Bool("federalonly", true, "only print federal holidays, use -federalonly=false to print all. See README")
	color := flag.Bool("color", true, "colorize holidays less than 30 days away, use -color=false to disable")
	savecsv := flag.Bool("savecsv", false, "saves holidays to 'holidays.csv' file, use -savecsv to enable")
	savejson := flag.Bool("savejson", false, "saves holidays to 'holidays.json' file, use -savejson to enable")
	stdout := flag.Bool("stdout", false, "output JSON to stdout and exit. See README.")
	flag.Parse()

	// Fetch country list from Nager API
	if *listCountries == true {
		countries, err := api.ListCountries(client, *debug)
		if err != nil {
			fmt.Printf("Error: %s.\n", err)
			return
		}

		// Send countries to stdout if stdout flag is sent and end program
		if *stdout {
			err := api.CountryStdout(countries)
			if err != nil {
				fmt.Printf("Error: %s.\n", err)
				return
			}
			return
		}

		// Print to terminal
		api.PrintCountries(countries)

		// Save to .csv
		if *savecsv {
			err := api.CountryOutputCSV(countries)
			if err != nil {
				fmt.Printf("Error: %s.\n", err)
				return
			}
		}

		// Save to .json
		if *savejson {
			err := api.CountryOutputJSON(countries)
			if err != nil {
				fmt.Printf("Error: %s.\n", err)
				return
			}
		}
		return
	}

	// Get holidays from Nager API
	rawHolidays, err := api.GetHolidays(client, *year, *countryCode, *debug)
	if err != nil {
		fmt.Printf("Error: %s.\n", err)
		return
	}

	// Enrich holidays
	holidays, err := api.EnrichHolidays(rawHolidays)

	// Send holidays to stdout if flag is sent and end program
	if *stdout {
		err := api.HolidayStdout(holidays)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		return
	}
	// Save to .csv
	if *savecsv {
		err := api.HolidayOutputCSV(holidays)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
	}
	// Save to .json
	if *savejson {
		err := api.HolidayOutputJSON(holidays)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
	}
	// Print holidays to terminal is not saving to files or piping to stdout
	if !*savecsv && !*savejson {
		api.PrintHolidays(holidays, *year, *countryCode, *federalOnly, *color)
	}
}
