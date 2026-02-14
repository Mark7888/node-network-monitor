package validators

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ValidateUUID checks if a string is a valid UUID
func ValidateUUID(s string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID format")
	}
	return id, nil
}

// ValidateTimeRange checks if a time range is valid
func ValidateTimeRange(from, to *time.Time) error {
	if from != nil && to != nil && from.After(*to) {
		return fmt.Errorf("from time must be before to time")
	}
	return nil
}

// ValidatePagination validates and normalizes pagination parameters
func ValidatePagination(page, limit int) (int, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}
	if limit > 1000 {
		limit = 1000
	}
	return page, limit, nil
}

// ValidateInterval checks if an interval string is valid
func ValidateInterval(interval string) error {
	validIntervals := map[string]bool{
		"1h": true,
		"6h": true,
		"1d": true,
	}
	if !validIntervals[interval] {
		return fmt.Errorf("invalid interval, must be one of: 1h, 6h, 1d")
	}
	return nil
}
