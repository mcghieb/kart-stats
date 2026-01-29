package scraper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/mcghieb/kart-stats/web-scraper/internal/models"
	"github.com/mcghieb/kart-stats/web-scraper/internal/scraper/cache"
)

func setDriverID(e *colly.HTMLElement, c *cache.Cache) error {
	// fetch alias
	alias := e.ChildText("td.Racername > span > a")

	// get driver from cache based on alias
	d, exists := c.Driver.GetByAlias(alias)
	if !exists {
		return setFromHTMLError(fmt.Sprintf("driver from [%s]", alias))
	}

	updateResultInCache(e, c, func(r *models.RaceResult) {
		r.DriverID = d.ID
	})
	return nil
}

func setBestLaptime(e *colly.HTMLElement, c *cache.Cache) error {
	// get best laptime
	blt := e.ChildText("td.BestLap > span")

	bestLapTime, err := strconv.ParseFloat(blt, 64)
	if err != nil {
		return setFromHTMLError("best lap time")
	}

	updateResultInCache(e, c, func(r *models.RaceResult) {
		r.BestLaptime = bestLapTime
	})
	return nil
}

func setNumLaps(e *colly.HTMLElement, c *cache.Cache) error {
	// get number of laps
	l := e.ChildText("td.Laps > span")

	numLaps, err := strconv.Atoi(l)
	if err != nil {
		return setFromHTMLError("number of laps")
	}

	updateResultInCache(e, c, func(r *models.RaceResult) {
		r.NumLaps = numLaps
	})
	return nil
}

func setLeaderGap(e *colly.HTMLElement, c *cache.Cache) error {
	// get gap from leader
	g := e.ChildText("td.Gap > span")

	var gap float64
	if g == "-" {
		gap = 0.00
	} else {
		var err error
		gap, err = strconv.ParseFloat(g, 64) // TODO: this might cause problems (test this)
		if err != nil {
			return setFromHTMLError("gap from leader")
		}
	}

	updateResultInCache(e, c, func(r *models.RaceResult) {
		r.GapFromLeader = gap
	})
	return nil
}

func setAvgLaptime(e *colly.HTMLElement, c *cache.Cache) error {
	// get average laptime
	al := e.ChildText("td.AvgLap > span")

	avgLaptime, err := strconv.ParseFloat(al, 64)
	if err != nil {
		return setFromHTMLError("average laptime")
	}

	updateResultInCache(e, c, func(r *models.RaceResult) {
		r.AvgLaptime = avgLaptime
	})
	return nil
}

func setProskill(e *colly.HTMLElement, c *cache.Cache) error {
	// get proskill rating
	pr := e.ChildText("td.RPM > span")

	proskillRating, err := strconv.Atoi(pr)
	if err != nil {
		// TODO: ERROR
		return setFromHTMLError("proskill rating")
	}

	updateResultInCache(e, c, func(r *models.RaceResult) {
		r.ProskillRating = proskillRating
	})
	driverID := c.Race.Get(getHeatNum(e)).DriverID
	if err := c.Driver.UpdateProskill(driverID, proskillRating); err != nil {
		return setFromHTMLError(fmt.Sprintf("proskill rating: %s", err))
	}
	return nil
}

func setPosition(e *colly.HTMLElement, c *cache.Cache) error {
	p := e.ChildText("td.Position > span")

	position, err := strconv.Atoi(p)
	if err != nil {
		return setFromHTMLError("position")
	}

	updateResultInCache(e, c, func(r *models.RaceResult) {
		r.Position = position
	})
	return nil
}

func getHeatNum(e *colly.HTMLElement) string {
	url := e.Request.URL.String()
	n := strings.Split(url, "=")[1]

	return n
}

func setFromHTMLError(s string) error {
	return fmt.Errorf("Error: failed to fetch %s\n", s)
}

func updateResultInCache(e *colly.HTMLElement, c *cache.Cache, updater func(*models.RaceResult)) {
	heatNum := getHeatNum(e)
	result := c.Race.Get(heatNum)
	updater(&result)
	c.Race.Put(result)
}
