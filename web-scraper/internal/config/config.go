package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL       string
	KartTrackDomainID string // this is what comes before ".clubspeed.com"
	StartHeatNo       int
	EndHeatNo         int
}

func LoadConfig() Config {
	startHeat := os.Getenv("START_NO")
	endHeat := os.Getenv("END_NO")

	startHeatNo, _ := strconv.Atoi(startHeat)
	endHeatNo, _ := strconv.Atoi(endHeat)

	return Config{
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		KartTrackDomainID: os.Getenv("KART_TRACK_DOMAIN_ID"),
		StartHeatNo:       startHeatNo,
		EndHeatNo:         endHeatNo,
	}
}
