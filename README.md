# Holiday Lookup

A command-line tool that fetches public holidays for a given country and year using the [Nager.Date](https://date.nager.at) public API.

> This project was built as a learning exercise while teaching myself Go.

## Features

- Look up public holidays for any supported country and year
- Defaults to the current year and the United States
- Shows the amount of days until next holiday
- Show only federally recognized (global) holidays, or all holidays
- List every country code the API supports
- Defaults to colorizing all holidays less than 30 days away to blue
- Debug mode to print the raw API response (for troubleshooting)

## Installation

```bash
git clone https://github.com/biggen1684/holidayapi.git
cd holidayapi
go build
```

## Usage

> **Note:** Windows users should run `holidayapi-windows-amd64.exe` from Command Prompt or PowerShell 
> rather than double-clicking the file, otherwise the terminal will close immediately 
> after the program exits.

Look up current year's US federal holidays (the default):

```bash
go run main.go
```

Look up holidays for a specific year and country:

```bash
go run main.go -year=2025 -countrycode=CA
```

Show all holidays, not just federal ones:

```bash
go run main.go -federalonly=false
```

> **Note:** Running with `-federalonly=false` may show some holidays more than once.
> The Nager.Date API returns separate entries for holidays that apply to different
> regions or have different observance types (e.g. public vs. bank holiday), so the
> same holiday can appear multiple times on the same date.

List all available country codes:

```bash
go run main.go -listcountries
```

Disable colorization:

```bash
go run main.go -color=false
```

Print the raw API response (useful for troubleshooting):

```bash
go run main.go -debug
```

## Flags

- `-year` — The year to look up, in `YYYY` format  
- `-countrycode` — Two-letter ISO 3166-1 alpha-2 country code  
- `-federalonly` — Show only federal holidays. Use `-federalonly=false` to show all  
- `-listcountries` — List all available country codes and exit  
- `-color` — Colorize holidays within 30 days in blue. Use `-color=false` to disable  
- `-debug` — Print the raw API response

## Example Output
> Holidays within 30 days are highlighted in blue in the terminal.
```
The holidays in 2026 for the country of US are as follows:

1. Thursday, 01-01 New Year's Day 
2. Monday, 01-19 Martin Luther King, Jr. Day 
3. Monday, 02-16 Presidents Day 
4. Monday, 05-25 Memorial Day 
5. Friday, 06-19 Juneteenth National Independence Day (14 days away)
6. Friday, 07-03 Independence Day (28 days away)
7. Monday, 09-07 Labour Day (94 days away)
8. Monday, 10-12 Columbus Day (129 days away)
9. Wednesday, 11-11 Veterans Day (159 days away)
10. Thursday, 11-26 Thanksgiving Day (174 days away)
11. Friday, 12-25 Christmas Day (203 days away)
```

## License

MIT
