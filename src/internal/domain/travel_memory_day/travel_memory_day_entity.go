package travelmemoryday

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidDayDate   = errors.New("day date must be within travel memory start and end dates")
	ErrDuplicateDayDate = errors.New("duplicate day date in travel memory")
)

type TravelMemoryDay struct {
	ID          int64
	MemoryID    int64
	Date        time.Time
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTravelMemoryDay(
	memoryID int64,
	date, memoryStartDate, memoryEndDate time.Time,
	description *string,
) (TravelMemoryDay, error) {
	day := TravelMemoryDay{MemoryID: memoryID}
	if err := day.Update(date, memoryStartDate, memoryEndDate, description); err != nil {
		return TravelMemoryDay{}, err
	}
	return day, nil
}

func (d *TravelMemoryDay) Update(
	date, memoryStartDate, memoryEndDate time.Time,
	description *string,
) error {
	normDate := normalizeDate(date)
	normStart := normalizeDate(memoryStartDate)
	normEnd := normalizeDate(memoryEndDate)

	if normDate.IsZero() || normStart.IsZero() || normEnd.IsZero() ||
		normDate.Before(normStart) || normDate.After(normEnd) {
		return ErrInvalidDayDate
	}

	d.Date = normDate
	d.Description = trimOptionalString(description)
	return nil
}

func normalizeDate(t time.Time) time.Time {
	if t.IsZero() {
		return time.Time{}
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func trimOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}
