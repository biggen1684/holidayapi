package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/biggen1684/holidayapi/api"
)

func main() {

	client := &http.Client{Timeout: 30 * time.Second}

	//Create flags to pass into CLI. Defaults to current year and "US" while debug disabled by default
	currentYear := fmt.Sprintf("%d", time.Now().Year())
	year := flag.String("year", currentYear, "the year in YYYY format")
	countryCode := flag.String("countrycode", "US", "2-letter ISO 3166-1 alpha-2 country code")
	debug := flag.Bool("debug", false, "print raw API response (use -debug to enable)")
	listCountries := flag.Bool("listcountries", false, "list all available countries and their 2 letter codes (use -listcountries to enable)")
	federalOnly := flag.Bool("federalonly", true, "only show federal holidays - may show duplicates (use -federalonly=false to show all)")
	color := flag.Bool("color", true, "colorize holidays less than 30 days away (use -color=false to disable)")
	savecsv := flag.Bool("savecsv", false, "saves holidays to 'holidays.csv' file (use -savecsv to enable)")
	//savejson := flag.Bool("savejson", false, "saves holidays to 'holidays.json' file (use -savejson to enable)")
	flag.Parse()

	//List countries if flag is passed in and then terminate program
	if *listCountries == true {
		countries, err := api.ListCountries(client, *debug)
		if err != nil {
			fmt.Printf("Error: %s.\n", err)
			os.Exit(1)
		}
		api.PrintCountries(countries)
		return
	}

	//Get holidays from Nager API
	rawHolidays, err := api.GetHolidays(client, *year, *countryCode, *debug)
	if err != nil {
		fmt.Printf("Error: %s.\n", err)
		return
	}

	//Enrich holidays
	holidays, err := api.EnrichHolidays(rawHolidays)

	//Output either .csv, .json, or terminal depending on flags used
	switch {
	case *savecsv:
		err := api.OutputCSV(holidays)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
	// case *savejson:
	// 	err := api.OutputJSON(holidays)
	// 	if err != nil {
	// 		fmt.Printf("Error: %s\n", err)
	// 		return
	// 	}
	default:
		api.PrintHolidays(holidays, *year, *countryCode, *federalOnly, *color)
	}
}
