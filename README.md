# Holiday Lookup

A command-line tool that fetches public holidays for a given country and year using the [Nager.Date](https://date.nager.at) public API.

> This project was built as a learning exercise while teaching myself Go.

## Features

- Look up public holidays for any supported country and year
- Defaults to the current year and the United States
- Show only federally recognized (global) holidays, or all holidays
- List every country code the API supports
- Debug mode to print the raw API response (for troubleshooting)

## Installation

```bash
git clone https://github.com/biggen1684/holidayapi.git
cd holidayapi
go build
```

## Usage

> **Note:** Windows users should run `holidayapi.exe` from Command Prompt or PowerShell 
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
go run main.go -globalonly=false
```

> **Note:** Running with `-globalonly=false` may show some holidays more than once.
> The Nager.Date API returns separate entries for holidays that apply to different
> regions or have different observance types (e.g. public vs. bank holiday), so the
> same holiday can appear multiple times on the same date.

List all available country codes:

```bash
go run main.go -listcountries
```

Print the raw API response (useful for troubleshooting):

```bash
go run main.go -debug
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-year` | current year | The year to look up, in `YYYY` format |
| `-countrycode` | `US` | Two-letter ISO 3166-1 alpha-2 country code |
| `-globalonly` | `true` | Show only federal/global holidays. Use `-globalonly=false` to show all |
| `-listcountries` | `false` | List all available country codes and exit |
| `-debug` | `false` | Print the raw API response |

## Example Output

```
The holidays in 2026 for country US are as follows:

1. 01-01 New Year's Day
2. 01-19 Martin Luther King, Jr. Day
3. 02-16 Presidents Day
4. 05-25 Memorial Day
5. 06-19 Juneteenth National Independence Day
6. 07-03 Independence Day
7. 09-07 Labour Day
8. 10-12 Columbus Day
9. 11-11 Veterans Day
10. 11-26 Thanksgiving Day
11. 12-25 Christmas Day
```

## License

MIT
