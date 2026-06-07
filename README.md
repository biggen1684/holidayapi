# Holiday Lookup

A command-line tool that fetches public holidays for a given country and year using the [Nager.Date](https://date.nager.at) public API.

> This project was built as a learning exercise while teaching myself Go.

## Features

- Look up public holidays for any supported country and year
- Defaults to the current year and the United States
- Prints the number of days until each upcoming holiday
- Prints only federally recognized holidays, or all holidays
- Prints if today is a holiday
- List every country code the API supports
- Colorizes holidays within 30 days in blue (enabled by default)
- Option to save output to a `.csv` file
- Option to save output to a `.json` file
- Debug mode to print the raw API response (for troubleshooting)

## Installation

**Option 1 — Download a pre-built binary** from the [Releases](https://github.com/biggen1684/holidayapi/releases) page and run it directly from a terminal.

**Option 2 — Build from source** (requires Go):

```bash
git clone https://github.com/biggen1684/holidayapi.git
cd holidayapi
go build -o holidayapi
./holidayapi
```

## Usage

> **Note:** Windows users should run `holidayapi-windows-amd64.exe` from Command Prompt or PowerShell
> rather than double-clicking the file, otherwise the terminal will close immediately
> after the program exits.

Look up current year's US federal holidays (the default):

```bash
./holidayapi
```

Look up holidays for a specific year and country:

```bash
./holidayapi -year=2025 -countrycode=CA
```

Show all holidays, not just federal ones:

```bash
./holidayapi -federalonly=false
```

> **Note:** Running with `-federalonly=false` may show some holidays more than once.
> The Nager.Date API returns separate entries for holidays that apply to different
> regions or have different observance types (e.g. public vs. bank holiday), so the
> same holiday can appear multiple times on the same date.

List all available country codes:

```bash
./holidayapi -listcountries
```

Save results to a CSV file:

```bash
./holidayapi -savecsv
```

Save results to a JSON file:

```bash
./holidayapi -savejson
```

Disable colorization:

```bash
./holidayapi -color=false
```

Print the raw API response (useful for troubleshooting):

```bash
./holidayapi -debug
```

## Flags

- `-year` — The year to look up, in `YYYY` format. Default: current year.
- `-countrycode` — Two-letter ISO 3166-1 alpha-2 country code. Default: `US`.
- `-federalonly` — Show only federal holidays. Default: `true`. Use `-federalonly=false` to show all.
- `-listcountries` — List all available country codes and exit.
- `-savecsv` — Save holidays to `holidays.csv`. Use `-savecsv` to enable.
- `-savejson` — Save holidays to `holidays.json`. Use `-savejson` to enable.
- `-color` — Colorize holidays within 30 days in blue. Default: `true`. Use `-color=false` to disable.
- `-debug` — Print the raw API response. Use `-debug` to enable.

## Example Output

> Holidays within 30 days are highlighted in blue in the terminal.

```
The holidays in 2026 for the country of US are as follows:

🎉 Today is Independence Day!

1. Thursday, 01-01 New Year's Day 
2. Monday, 01-19 Martin Luther King, Jr. Day 
3. Monday, 02-16 Presidents Day 
4. Monday, 05-25 Memorial Day 
5. Friday, 06-19 Juneteenth National Independence Day (12 days away)
6. Friday, 07-03 Independence Day (26 days away)
7. Monday, 09-07 Labour Day (92 days away)
8. Monday, 10-12 Columbus Day (127 days away)
9. Wednesday, 11-11 Veterans Day (157 days away)
10. Thursday, 11-26 Thanksgiving Day (172 days away)
11. Friday, 12-25 Christmas Day (201 days away)
```

## Output Fields for JSON and CSV

| Field | Type | Description |
|-------|------|-------------|
| `date` | string | Holiday date in `YYYY-MM-DD` format |
| `countryCode` | string | Two-letter ISO country code |
| `name` | string | Name of the holiday |
| `global` | bool | `true` if the holiday is observed nationwide |
| `weekday` | string | Day of the week the holiday falls on |
| `underThirty` | bool | `true` if the holiday is less than 30 days from today |
| `daysAway` | int | Days until the holiday. Negative values indicate days elapsed since the holiday passed |
| `isToday` | bool | `true` if holiday is today |

## CSV Output Example

```

date,countryCode,name,global,weekday,underThirty,daysAway,isToday
2026-01-01,US,New Year's Day,true,Thursday,false,-157,false
2026-01-19,US,"Martin Luther King, Jr. Day",true,Monday,false,-139,false
2026-02-12,US,Lincoln's Birthday,false,Thursday,false,-115,false
2026-02-16,US,Presidents Day,true,Monday,false,-111,false
2026-04-03,US,Good Friday,false,Friday,false,-65,false
2026-04-03,US,Good Friday,false,Friday,false,-65,false
2026-05-08,US,Truman Day,false,Friday,false,-30,false
2026-05-25,US,Memorial Day,true,Monday,false,-13,false
2026-06-19,US,Juneteenth National Independence Day,true,Friday,true,12,false
2026-07-03,US,Independence Day,true,Friday,true,26,true
2026-09-07,US,Labour Day,true,Monday,false,92,false
2026-10-12,US,Columbus Day,false,Monday,false,127,false
2026-10-12,US,Columbus Day,true,Monday,false,127,false
2026-10-12,US,Indigenous Peoples' Day,false,Monday,false,127,false
2026-11-11,US,Veterans Day,true,Wednesday,false,157,false
2026-11-26,US,Thanksgiving Day,true,Thursday,false,172,false
2026-12-25,US,Christmas Day,true,Friday,false,201,false
```

## JSON Output Example

[{"date":"2026-01-01","countryCode":"US","name":"New Year's Day","global":true,"weekday":"Thursday","underThirty":false,"daysAway":-157,"isToday":false},{"date":"2026-01-19","countryCode":"US","name":"Martin Luther King, Jr. Day","global":true,"weekday":"Monday","underThirty":false,"daysAway":-139,"isToday":false},{"date":"2026-02-12","countryCode":"US","name":"Lincoln's Birthday","global":false,"weekday":"Thursday","underThirty":false,"daysAway":-115,"isToday":false},{"date":"2026-02-16","countryCode":"US","name":"Presidents Day","global":true,"weekday":"Monday","underThirty":false,"daysAway":-111,"isToday":false},{"date":"2026-04-03","countryCode":"US","name":"Good Friday","global":false,"weekday":"Friday","underThirty":false,"daysAway":-65,"isToday":false}...]

## License

MIT
