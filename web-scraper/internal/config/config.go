package config

import "os"

type Config struct {
	DatabaseURL       string
	KartTrackDomainID string // this is what comes before ".clubspeed.com"
}

func LoadConfig() Config {
	return Config{
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		KartTrackDomainID: os.Getenv("KART_TRACK_DOMAIN_ID"),
	}
}
