package driverhistory

import (
	"fmt"
	"log"

	"github.com/gocolly/colly/v2"
	"github.com/mcghieb/kart-stats/web-scraper/internal/cache"
)

// AttachHandlers attaches handlers for scraping driver history pages
// This scraper is used to get current proskill ratings for all drivers
func AttachHandlers(c *colly.Collector, appCache *cache.Driver) {
	// Parse current proskill rating
	c.OnHTML("span#lblSpeedLimit", func(e *colly.HTMLElement) {
		driverID := e.Request.Ctx.Get("driverID")
		if driverID == "" {
			log.Printf("missing driverID in context")
			return
		}

		fmt.Printf("PROSKILL for %s is %s\n", driverID, e.Text)
		proskill, err := parseProskill(e.Text)
		if err != nil {
			log.Printf("failed to parse proskill for driver %s: %v", driverID, err)
			return
		}

		appCache.UpdateProskill(driverID, proskill)
	})
}
