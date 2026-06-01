package main

import (
	"flag"
	"fmt"
	"holiday/api"
	"net/http"
	"os"
	"time"
)

func main() {

	client := &http.Client{Timeout: 30 * time.Second}

	//Create flags to pass into CLI. Defaults to current year and "US" while debug disabled by default
	currentYear := fmt.Sprintf("%d", time.Now().Year())
	year := flag.String("year", currentYear, "the year in xxxx format")
	countryCode := flag.String("countrycode", "US", "2-letter ISO 3166-1 alpha-2 country code")
	debug := flag.Bool("debug", false, "print raw API response (use -debug to enable)")
	listCountries := flag.Bool("listcountries", false, "list all available countries (use -listcountries to enable)")
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

	holidays, err := api.GetHolidays(client, *year, *countryCode, *debug)
	if err != nil {
		fmt.Printf("Error: %s.\n", err)
		return
	}
	api.PrintHolidays(holidays, *year, *countryCode)

}
