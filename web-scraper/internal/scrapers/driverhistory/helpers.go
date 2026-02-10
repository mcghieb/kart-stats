package driverhistory

import (
	"fmt"
	"strconv"
)

func parseProskill(text string) (int, error) {
	proskill, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("failed to parse proskill from [%s]: %w", text, err)
	}
	return proskill, nil
}
