package travelmemory

import (
	"errors"
	"path"
	"strings"
	"time"
)

var (
	ErrInvalidDateRange = errors.New("end date must be on or after start date")
	ErrInvalidImagePath = errors.New("image path must start with travel_memories/")
)

type TravelMemory struct {
	ID             int64
	LocationID     int64
	Title          string
	StartDate      time.Time
	EndDate        time.Time
	Description    *string
	CoverImagePath string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewTravelMemory(
	locationID int64,
	title string,
	startDate, endDate time.Time,
	description, coverImagePath *string,
) (TravelMemory, error) {
	memory := TravelMemory{}

	if err := memory.Update(
		locationID,
		title,
		startDate,
		endDate,
		description,
		coverImagePath,
	); err != nil {
		return TravelMemory{}, err
	}

	return memory, nil
}

func (m *TravelMemory) Update(
	locationID int64,
	title string,
	startDate, endDate time.Time,
	description, coverImagePath *string,
) error {
	normStart := normalizeDate(startDate)
	normEnd := normalizeDate(endDate)

	if normStart.IsZero() || normEnd.IsZero() || normEnd.Before(normStart) {
		return ErrInvalidDateRange
	}

	var normalizedCoverPath string
	trimmedCover := trimOptionalString(coverImagePath)
	if trimmedCover != nil {
		if !isValidImagePath(*trimmedCover) {
			return ErrInvalidImagePath
		}
		normalizedCoverPath = *trimmedCover
	}

	m.LocationID = locationID
	m.Title = strings.TrimSpace(title)
	m.StartDate = normStart
	m.EndDate = normEnd
	m.Description = trimOptionalString(description)
	m.CoverImagePath = normalizedCoverPath

	return nil
}

func normalizeDate(t time.Time) time.Time {
	if t.IsZero() {
		return time.Time{}
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func isValidImagePath(value string) bool {
	const prefix = "travel_memories/"

	if !strings.HasPrefix(value, prefix) {
		return false
	}

	relativePath := strings.TrimPrefix(value, prefix)

	if relativePath == "" ||
		strings.Contains(relativePath, "\\") ||
		path.IsAbs(relativePath) ||
		path.Clean(relativePath) != relativePath ||
		relativePath == ".." ||
		strings.HasPrefix(relativePath, "../") {
		return false
	}

	return true
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
