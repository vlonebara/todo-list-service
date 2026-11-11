package service

import (
	"fmt"
	"time"
)

const dateLayout = "20060102"

// ParseDate проверяет общий формат и поддерживаемый диапазон календарной даты.
func ParseDate(value string) (time.Time, error) {
	if len(value) != 8 {
		return time.Time{}, fmt.Errorf("date must use YYYYMMDD format")
	}
	date, err := time.Parse(dateLayout, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date: %w", err)
	}
	if date.Year() < 1 {
		return time.Time{}, fmt.Errorf("date year must be between 1 and 9999")
	}
	return date, nil
}
