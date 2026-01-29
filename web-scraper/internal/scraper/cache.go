package scraper

import (
	"fmt"
	"sync"

	"github.com/mcghieb/kart-stats/web-scraper/internal/models"
)

type DriverCache struct {
	drivers sync.Map
}

func NewDriverCache() *DriverCache {
	return &DriverCache{}
}

func (c *DriverCache) Put(d models.Driver) {
	c.drivers.Store(d.ID, d)
}

func (c *DriverCache) Has(id string) bool {
	_, exists := c.drivers.Load(id)
	return exists
}

func (c *DriverCache) UpdateProskill(id string, p int) error {
	val, ok := c.drivers.Load(id)
	if !ok {
		return fmt.Errorf(
			"ERROR: attempted to update cache object [%s] that doesn't exist",
			id,
		)
	}

	d := val.(models.Driver)
	d.ProskillRating = p
	c.Put(d)

	return nil
}
