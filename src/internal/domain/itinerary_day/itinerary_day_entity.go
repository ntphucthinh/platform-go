package itineraryday

import (
	"errors"
	"path"
	"strings"
	"time"
)

var (
	ErrDayNotFound        = errors.New("itinerary day not found")
	ErrInvalidDayNumber   = errors.New("day number must be greater than zero")
	ErrDuplicateDayNumber = errors.New("duplicate day number in itinerary")
	ErrInvalidImagePath   = errors.New("image path must start with itineraries/days/")
)

type ItineraryDay struct {
	ID          int64
	ItineraryID int64
	DayNumber   int
	Title       string
	Subtitle    *string
	Summary     *string
	ImagePath   *string
	DayTips     []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewItineraryDay(
	itineraryID int64,
	dayNumber int,
	title string,
	subtitle, summary, imagePath *string,
	dayTips []string,
) (ItineraryDay, error) {
	day := ItineraryDay{ItineraryID: itineraryID}
	if err := day.Update(dayNumber, title, subtitle, summary, imagePath, dayTips); err != nil {
		return ItineraryDay{}, err
	}
	return day, nil
}

func (d *ItineraryDay) Update(
	dayNumber int,
	title string,
	subtitle, summary, imagePath *string,
	dayTips []string,
) error {
	if dayNumber <= 0 {
		return ErrInvalidDayNumber
	}

	trimmedImage := trimOptionalString(imagePath)
	if trimmedImage != nil {
		if !isValidImagePath(*trimmedImage) {
			return ErrInvalidImagePath
		}
	}

	d.DayNumber = dayNumber
	d.Title = strings.TrimSpace(title)
	d.Subtitle = trimOptionalString(subtitle)
	d.Summary = trimOptionalString(summary)
	d.ImagePath = trimmedImage
	d.DayTips = normalizeStringSlice(dayTips)

	return nil
}

func isValidImagePath(value string) bool {
	const prefix = "itineraries/days/"

	if !strings.HasPrefix(value, prefix) {
		return false
	}

	relativePath := strings.TrimPrefix(value, prefix)

	if relativePath == "" ||
		strings.Contains(relativePath, "\\") ||
		path.IsAbs(relativePath) ||
		path.Clean(relativePath) != relativePath ||
		relativePath == "." ||
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

func normalizeStringSlice(items []string) []string {
	if len(items) == 0 {
		return nil
	}

	normalized := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			normalized = append(normalized, trimmed)
		}
	}
	if len(normalized) == 0 {
		return nil
	}
	return normalized
}
