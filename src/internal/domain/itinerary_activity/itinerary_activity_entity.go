package itineraryactivity

import (
	"errors"
	"math"
	"path"
	"strings"
	"time"
)

var (
	ErrActivityNotFound       = errors.New("itinerary activity not found")
	ErrInvalidStartTime       = errors.New("start time must be in HH:mm format")
	ErrInvalidDurationMinutes = errors.New("duration minutes must be between 1 and 65535")
	ErrInvalidImagePath       = errors.New("image path must start with itineraries/activities/")
)

type ItineraryActivity struct {
	ID               int64
	ItineraryDayID   int64
	StartTime        string
	Title            string
	ActivityCategory *string
	DurationMinutes  *uint16
	Location         *string
	GoogleMapsURL    *string
	Description      *string
	Notes            *string
	ImagePath        *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func NewItineraryActivity(
	itineraryDayID int64,
	startTime, title string,
	activityCategory *string,
	durationMinutes *int,
	location, googleMapsURL, description, notes, imagePath *string,
) (ItineraryActivity, error) {
	act := ItineraryActivity{ItineraryDayID: itineraryDayID}
	if err := act.Update(
		startTime,
		title,
		activityCategory,
		durationMinutes,
		location,
		googleMapsURL,
		description,
		notes,
		imagePath,
	); err != nil {
		return ItineraryActivity{}, err
	}
	return act, nil
}

func (a *ItineraryActivity) Update(
	startTime, title string,
	activityCategory *string,
	durationMinutes *int,
	location, googleMapsURL, description, notes, imagePath *string,
) error {
	normalizedTime, err := normalizeStartTime(startTime)
	if err != nil {
		return err
	}

	var normDuration *uint16
	if durationMinutes != nil {
		if *durationMinutes <= 0 || *durationMinutes > math.MaxUint16 {
			return ErrInvalidDurationMinutes
		}
		d := uint16(*durationMinutes)
		normDuration = &d
	}

	trimmedImage := trimOptionalString(imagePath)
	if trimmedImage != nil {
		if !isValidImagePath(*trimmedImage) {
			return ErrInvalidImagePath
		}
	}

	a.StartTime = normalizedTime
	a.Title = strings.TrimSpace(title)
	a.ActivityCategory = trimOptionalString(activityCategory)
	a.DurationMinutes = normDuration
	a.Location = trimOptionalString(location)
	a.GoogleMapsURL = trimOptionalString(googleMapsURL)
	a.Description = trimOptionalString(description)
	a.Notes = trimOptionalString(notes)
	a.ImagePath = trimmedImage

	return nil
}

func normalizeStartTime(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if t, err := time.Parse("15:04", trimmed); err == nil && len(trimmed) == 5 {
		return t.Format("15:04:05"), nil
	}
	if t, err := time.Parse("15:04:05", trimmed); err == nil && len(trimmed) == 8 {
		return t.Format("15:04:05"), nil
	}
	return "", ErrInvalidStartTime
}

func isValidImagePath(value string) bool {
	const prefix = "itineraries/activities/"

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
