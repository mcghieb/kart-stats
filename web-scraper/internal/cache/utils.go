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

func UpdateCachedDriverID(
	e *colly.HTMLElement,
	c *Cache,
	heatNum string,
	aliasParser func(*colly.HTMLElement) (string, error),
) error {
	alias, err := aliasParser(e)
	if err != nil {
		return err
	}

	driver, exists := c.Driver.ByAlias(alias)
	if !exists {
		return fmt.Errorf("ERROR: attempting to get a driver that does not exist from the cache using driver alias")
	}

	UpdateCachedRace(c, heatNum, func(r *models.Race) {
		r.DriverID = driver.ID
	})

	return nil
}
