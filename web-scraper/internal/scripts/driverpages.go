package scripts

import (
	"fmt"
	"log"

	"github.com/gocolly/colly/v2"
	"github.com/mcghieb/kart-stats/web-scraper/internal/cache"
	"github.com/mcghieb/kart-stats/web-scraper/internal/collector"
	"github.com/mcghieb/kart-stats/web-scraper/internal/scrapers/driverhistory"
)

const (
	BASEURL = "https://rrorem.clubspeed.com/sp_center/RacerHistory.aspx?CustID"
)

func ParseDriverPages(appCache *cache.Cache) {
	c := collector.NewCollector()
	driverhistory.AttachHandlers(c, appCache)

	// Build list of URLs from all drivers in cache
	type visit struct {
		url      string
		driverID string
	}

	var visits []visit
	appCache.Driver.Range(func(id string) bool {
		visits = append(visits, visit{
			url:      createURL(id),
			driverID: id,
		})
		return true
	})

	fmt.Printf("Visiting %d driver history pages\n", len(visits))

	for _, v := range visits {
		ctx := colly.NewContext()
		ctx.Put("driverID", v.driverID)

		if err := c.Request("GET", v.url, nil, ctx, nil); err != nil {
			log.Printf("failed to visit driver page for %s: %v", v.driverID, err)
		}
	}
}

func createURL(driverID string) string {
	return fmt.Sprintf("%s=%s", BASEURL, driverID)
}
