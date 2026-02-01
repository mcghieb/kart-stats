package parse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

// HeatNum parses a heat number from the Race Heat Result specific page
func HeatNum(e *colly.HTMLElement) string {
	url := e.Request.URL.String()
	return strings.Split(url, "=")[1]
}

// BestLaptime parses the driver's best laptime from the Race Heat Result specific page
func BestLaptime(e *colly.HTMLElement) (float64, error) {
	return parseFloat(e, "td.BestLap > span", "best lap time", false)
}

// NumLaps parses the number of laps a driver took from the Race Heat Result specific page
func NumLaps(e *colly.HTMLElement) (int, error) {
	return parseInt(e, "td.Laps > span", "number of laps")
}

// LeaderGap parses the laptime gap from the leader of the race from the Race Heat Result specific page
func LeaderGap(e *colly.HTMLElement) (float64, error) {
	return parseFloat(e, "td.Gap > span", "gap from leader", true)
}

// AvgLaptime parses a driver's average laptime from the Race Heat Result specific page
func AvgLaptime(e *colly.HTMLElement) (float64, error) {
	return parseFloat(e, "td.AvgLap > span", "average laptime", false)
}

// Position parses a driver's position in the race from the Race Heat Result specific page
func Position(e *colly.HTMLElement) (int, error) {
	return parseInt(e, "td.Position > span", "position")
}

// Top3Position parses a top 3 driver's position in the race from the Race Heat Result specific page
func Top3Position(e *colly.HTMLElement) (int, error) {
	p := e.ChildText("td.Position") // this is a special parse specific to the top rows
	position := 0
	switch p {
	case "Heat Winner:":
		position = 1
	case "2nd Place:":
		position = 2
	case "3rd Place:":
		position = 3
	default:
		return position, generateFetchError("top 3 row postion")
	}

	return position, nil
}

// DriverAlias parses a driver's alias from the Race Heat Result specific page
func DriverAlias(e *colly.HTMLElement) (string, error) {
	return parseString(e, "td.Racername > span > a", "driver alias")
}

// Proskill parses a driver's proskill rating from the Race Heat Result specific page
func Proskill(e *colly.HTMLElement) (int, error) {
	return parseInt(e, "td.RPM > span", "proskill rating")
}

// generateFetchError represents a generator for an error
// when fetching data from an html page
func generateFetchError(s string) error {
	return fmt.Errorf("Error: failed to fetch %s\n", s)
}

// parseString takes an html element and a selector
// and parses for a string based on the selector
func parseString(
	e *colly.HTMLElement,
	selector,
	category string,
) (string, error) {
	val := e.ChildText(selector)
	if val == "" {
		return "", generateFetchError(category)
	}

	return val, nil
}

// parseFloat wraps parseNumber() to return an float from a specific css selector
func parseFloat(
	e *colly.HTMLElement,
	selector,
	category string,
	checkLeaderGap bool,
) (float64, error) {
	return parseNumber(e, selector, category, func(s string) (float64, error) {
		if checkLeaderGap && s == "-" { // if the driver is the leader for leader gap parsing
			return 0.00, nil
		}
		return strconv.ParseFloat(s, 64)
	})
}

// parseInt wraps parseNumber() to return an int from a specific css selector
func parseInt(
	e *colly.HTMLElement,
	selector,
	category string,
) (int, error) {
	return parseNumber(e, selector, category, func(s string) (int, error) {
		return strconv.Atoi(s)
	})
}

// parseNumber is a generic template function
// that wraps parseString() and returns either
// a float64 or an int
func parseNumber[T int | float64](
	e *colly.HTMLElement,
	selector,
	category string,
	convert func(s string) (T, error),
) (T, error) {
	s, err := parseString(e, selector, category)
	if err != nil {
		return 0.0, generateFetchError(category)
	}

	return convert(s)
}
