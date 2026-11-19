package itinerarytheme

import (
	"errors"
	"time"
)

var (
	ErrDuplicateThemeID = errors.New("duplicate theme ID in itinerary")
	ErrInvalidThemeID   = errors.New("one or more travel theme IDs are invalid")
)

type ItineraryTheme struct {
	ID          int64
	ItineraryID int64
	ThemeID     int64
	CreatedAt   time.Time
}

func NewItineraryTheme(itineraryID, themeID int64) ItineraryTheme {
	return ItineraryTheme{
		ItineraryID: itineraryID,
		ThemeID:     themeID,
	}
}
