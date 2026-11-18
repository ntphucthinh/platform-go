package qsdto

// ItineraryDto represents an itinerary item returned by the query service.
type ItineraryDto struct {
	ID             int64   `json:"id" gorm:"column:id"`
	LocationID     *int64  `json:"locationId" gorm:"column:location_id"`
	LocationName   *string `json:"locationName" gorm:"column:location_name"`
	ItineraryName  string  `json:"itineraryName" gorm:"column:itinerary_name"`
	DurationDays   int     `json:"durationDays" gorm:"-"`
	DurationNights int     `json:"durationNights" gorm:"-"`
	Duration       string  `json:"duration" gorm:"-"`
	ActivityCount  int     `json:"activityCount" gorm:"-"`
	Description    *string `json:"description" gorm:"column:description"`
	CoverImagePath *string `json:"coverImagePath" gorm:"column:cover_image_path"`
}
