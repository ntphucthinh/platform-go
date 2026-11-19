package itinerary

import (
	"errors"
	"path"
	"strings"
	"time"
)

var (
	ErrItineraryNotFound = errors.New("itinerary not found")
	ErrInvalidImagePath  = errors.New("image path must start with itineraries/")
)

type Itinerary struct {
	ID             int64
	ItineraryName  string
	LocationID     *int64
	SuitableFor    *string
	Description    *string
	CoverImagePath *string
	TravelTips     []string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewItinerary(
	itineraryName string,
	locationID *int64,
	suitableFor, description, coverImagePath *string,
	travelTips []string,
) (Itinerary, error) {
	it := Itinerary{}
	if err := it.Update(
		itineraryName,
		locationID,
		suitableFor,
		description,
		coverImagePath,
		travelTips,
	); err != nil {
		return Itinerary{}, err
	}
	return it, nil
}

func (it *Itinerary) Update(
	itineraryName string,
	locationID *int64,
	suitableFor, description, coverImagePath *string,
	travelTips []string,
) error {
	trimmedCover := trimOptionalString(coverImagePath)
	if trimmedCover != nil {
		if !isValidImagePath(*trimmedCover) {
			return ErrInvalidImagePath
		}
	}

	it.ItineraryName = strings.TrimSpace(itineraryName)
	it.LocationID = locationID
	it.SuitableFor = trimOptionalString(suitableFor)
	it.Description = trimOptionalString(description)
	it.CoverImagePath = trimmedCover
	it.TravelTips = normalizeStringSlice(travelTips)

	return nil
}

func isValidImagePath(value string) bool {
	const prefix = "itineraries/"

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
