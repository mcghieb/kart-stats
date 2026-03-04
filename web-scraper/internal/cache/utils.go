package cache

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/mcghieb/kart-stats/web-scraper/internal/models"
)

func UpdateCachedRaceAttribute[T any](e *colly.HTMLElement, c *Cache, heatNum string,
	parser func(*colly.HTMLElement) (T, error),
	updater func(*models.Race, T),
) error {
	value, err := parser(e)
	if err != nil {
		return err
	}
	UpdateCachedRace(c, heatNum, func(r *models.Race) {
		updater(r, value)
	})
	return nil
}

func UpdateCachedRace(c *Cache, heatNum string, setter func(*models.Race)) {
	race := c.Race.Get(heatNum) // get race
	setter(&race)               // set race attribute
	c.Race.Put(race)            // put race in cache
}

func UpdateCachedRaceLaps(c *Cache, heatNum string, driverID string, lapNum int, laptime float64) {
	race := c.Race.Get(heatNum)

	result := race.Results[driverID]
	result.Laps = append(result.Laps, models.Lap{LapNumber: lapNum, LapTime: laptime})
	race.Results[driverID] = result

	c.Race.Put(race)
}

// UpdateCachedResultAttribute updates a specific attribute of a driver's result in the race
func UpdateCachedResultAttribute[T any](
	e *colly.HTMLElement,
	c *Cache,
	heatNum string,
	driverID string,
	parser func(*colly.HTMLElement) (T, error),
	updater func(*models.Result, T),
) error {
	value, err := parser(e)
	if err != nil {
		return err
	}

	UpdateCachedRace(c, heatNum, func(r *models.Race) {
		if r.Results == nil {
			r.Results = make(map[string]models.Result)
		}
		result := r.Results[driverID]
		updater(&result, value)
		r.Results[driverID] = result
	})

	return nil
}

// UpdateCachedDriverID updates a result's driver ID for the current racer being processed
func UpdateCachedDriverID(
	e *colly.HTMLElement,
	c *Cache,
	heatNum string,
	aliasParser func(*colly.HTMLElement) (string, error),
) (string, error) {
	alias, err := aliasParser(e)
	if err != nil {
		return "", err
	}

	driver, exists := c.Driver.ByAlias(alias)
	if !exists {
		return "", fmt.Errorf("driver with alias %s not found in cache", alias)
	}

	UpdateCachedRace(c, heatNum, func(r *models.Race) {
		if r.Results == nil {
			r.Results = make(map[string]models.Result)
		}
		result := r.Results[driver.ID]
		result.DriverID = driver.ID
		r.Results[driver.ID] = result
	})

	return driver.ID, nil
}
