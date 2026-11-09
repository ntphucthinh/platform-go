package travellocation

import (
	"errors"
	"math"
	"path"
	"strings"
	"time"

	"github.com/hunter-io/geonames"
)

var (
	ErrInvalidCountryCode = errors.New("invalid country code")
	ErrInvalidCoordinates = errors.New("invalid coordinates")
	ErrInvalidImagePath   = errors.New("image path must start with travel_location/")
)

type TravelLocation struct {
	ID            int64
	CountryCode   string
	CountryName   string
	LocationName  string
	Latitude      *float64
	Longitude     *float64
	ImagePath     *string
	Status        string
	Description   *string
	Notes         *string
	GoogleMapsURL *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewTravelLocation(
	countryCode, countryName, locationName, status string,
	latitude, longitude *float64,
	imagePath, description, notes, googleMapsURL *string,
) (TravelLocation, error) {
	loc := TravelLocation{}

	if err := loc.Update(
		countryCode,
		countryName,
		locationName,
		status,
		latitude,
		longitude,
		imagePath,
		description,
		notes,
		googleMapsURL,
	); err != nil {
		return TravelLocation{}, err
	}

	return loc, nil
}

func (l *TravelLocation) Update(
	countryCode, countryName, locationName, status string,
	latitude, longitude *float64,
	imagePath, description, notes, googleMapsURL *string,
) error {
	normalizedCountry := strings.ToUpper(strings.TrimSpace(countryCode))

	if _, validCountry := geonames.SubdivisionsOf(normalizedCountry); !validCountry {
		return ErrInvalidCountryCode
	}

	if err := validateCoordinates(latitude, longitude); err != nil {
		return err
	}

	trimmedImagePath := trimOptionalString(imagePath)

	if trimmedImagePath != nil {
		if !isValidImagePath(*trimmedImagePath) {
			return ErrInvalidImagePath
		}
	}

	l.CountryCode = normalizedCountry
	l.CountryName = strings.TrimSpace(countryName)
	l.LocationName = strings.TrimSpace(locationName)
	l.Status = strings.TrimSpace(status)
	l.Latitude = latitude
	l.Longitude = longitude
	l.ImagePath = trimmedImagePath
	l.Description = trimOptionalString(description)
	l.Notes = trimOptionalString(notes)
	l.GoogleMapsURL = trimOptionalString(googleMapsURL)

	return nil
}

func validateCoordinates(latitude, longitude *float64) error {
	// Both coordinates must be provided together.
	if (latitude == nil) != (longitude == nil) {
		return ErrInvalidCoordinates
	}

	// Coordinates are optional.
	if latitude == nil {
		return nil
	}

	// Reject NaN and infinite values.
	if math.IsNaN(*latitude) || math.IsInf(*latitude, 0) ||
		math.IsNaN(*longitude) || math.IsInf(*longitude, 0) {
		return ErrInvalidCoordinates
	}

	// Validate latitude and longitude ranges.
	if *latitude < -90 || *latitude > 90 ||
		*longitude < -180 || *longitude > 180 {
		return ErrInvalidCoordinates
	}

	return nil
}

func isValidImagePath(value string) bool {
	const prefix = "travel_location/"

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
