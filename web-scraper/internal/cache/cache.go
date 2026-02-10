package cache

type Cache struct {
	Driver *Driver
	Race   *Race
}

func NewCache() *Cache {
	return &Cache{
		NewDriverCache(),
		NewRaceCache(),
	}
}
