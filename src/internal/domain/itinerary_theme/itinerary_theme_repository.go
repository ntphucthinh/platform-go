package itinerarytheme

import "context"

// ItineraryThemeRepository defines persistence operations for itinerary themes.
type ItineraryThemeRepository interface {
	// CountValidThemesByIDs returns the count of existing, non-deleted travel themes matching the given IDs.
	CountValidThemesByIDs(ctx context.Context, themeIDs []int64) (int64, error)
	// CreateItineraryThemes batch-inserts itinerary theme records.
	CreateItineraryThemes(ctx context.Context, themes []ItineraryTheme) error
	// DeleteThemesByItineraryID deletes all theme associations for a given itinerary ID.
	DeleteThemesByItineraryID(ctx context.Context, itineraryID int64) error
}
