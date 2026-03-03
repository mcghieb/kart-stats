package collector

import (
	"net/http"

	"github.com/gocolly/colly/v2"
)

// NewCollector is a wrapper for configuring colly.Collector
func NewCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains("rrorem.clubspeed.com"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:146.0) Gecko/20100101 Firefox/146.0"),
		colly.AllowURLRevisit(),
	)

	return c
}

// NewLocalCollector creates a collector for parsing local HTML files.
func NewLocalCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	c.WithTransport(t)

	return c
}
