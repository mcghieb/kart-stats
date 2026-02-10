package cache

import (
	"fmt"
	"sync"

	"github.com/mcghieb/kart-stats/web-scraper/internal/models"
)

type Driver struct {
	drivers   sync.Map
	aliasToID sync.Map
}

func NewDriverCache() *Driver {
	return &Driver{}
}

func (c *Driver) Get(id string) (models.Driver, bool) {
	val, exists := c.drivers.Load(id)
	if !exists {
		return models.Driver{}, false
	}

	d := val.(models.Driver)
	return d, true
}

func (c *Driver) ByAlias(a string) (models.Driver, bool) {
	id, exists := c.aliasToID.Load(a)
	if !exists {
		return models.Driver{}, false
	}

	return c.Get(id.(string))
}

func (c *Driver) Put(d models.Driver) {
	c.drivers.Store(d.ID, d)
	c.aliasToID.Store(d.Alias, d.ID)
}

func (c *Driver) Has(id string) bool {
	_, exists := c.drivers.Load(id)
	return exists
}

func (c *Driver) UpdateProskill(id string, p int) error {
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

func (c *Driver) String() string {
	result := "Driver Cache {\n"
	count := 0
	
	c.drivers.Range(func(key, value interface{}) bool {
		d := value.(models.Driver)
		result += fmt.Sprintf("  %s\n", d.String())
		count++
		return true
	})
	
	result += fmt.Sprintf("}\nTotal Drivers: %d", count)
	return result
}
