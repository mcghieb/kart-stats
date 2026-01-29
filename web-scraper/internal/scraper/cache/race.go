package cache

import "sync"

type Race struct {
	races sync.Map
}

func NewRaceCache() *Race {
	return &Race{}
}
