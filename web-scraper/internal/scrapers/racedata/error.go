package racedata

import "log"

// logErrors logs all errors in the slice using the standard logger
func logErrors(errs []error) {
	for _, err := range errs {
		if err != nil {
			log.Printf("scraper error: %v", err)
		}
	}
}
