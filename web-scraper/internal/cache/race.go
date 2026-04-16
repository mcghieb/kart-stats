package cache

import (
	"sync"

	"github.com/mcghieb/kart-stats/web-scraper/internal/models"
)

type Race struct {
	races sync.Map
}

func NewRaceCache() *Race {
	return &Race{}
}

func (r *Race) Put(result models.Race) {
	r.races.Store(result.ID, result)
}

func (r *Race) Get(raceID string) models.Race {
	val, exists := r.races.Load(raceID)
	if !exists {
		result := models.Race{ID: raceID}
		result.Results = make(map[string]models.Result)
		return result
	}

	return val.(models.Race)
}

func (r *Race) Has(raceID string) bool {
	_, exists := r.races.Load(raceID)
	return exists
}

func (r *Race) Range(f func(raceID string) bool) {
	r.races.Range(func(key, value interface{}) bool {
		return f(key.(string))
	})
}
