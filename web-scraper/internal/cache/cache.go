package cache

type Cache struct {
	Driver *Driver
	Race   *Race // FIXME: make this concurrency safe with threads...
}

func NewCache() *Cache {
	return &Cache{
		NewDriverCache(),
		NewRaceCache(),
	}
}
