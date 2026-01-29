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

func (r *Race) Put(result models.RaceResult) {
	r.races.Store(result.ID, result)
}

func (r *Race) Get(raceID string) models.RaceResult {
	val, exists := r.races.Load(raceID)
	if !exists {
		// put empty in cache, return empty
		result := models.RaceResult{ID: raceID}
		r.Put(result)
		return result
	}

	return val.(models.RaceResult)
}
