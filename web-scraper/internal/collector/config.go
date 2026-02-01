package collector

import "github.com/gocolly/colly/v2"

// NewCollector is a wrapper for configuring colly.Collector
func NewCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains("rrorem.clubspeed.com"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:146.0) Gecko/20100101 Firefox/146.0"),
	)

	return c
}
