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

// ItineraryDetailDto represents a detailed itinerary returned by the query service.
type ItineraryDetailDto struct {
	ID             int64                       `json:"id" gorm:"column:id"`
	ItineraryName  string                      `json:"itineraryName" gorm:"column:itinerary_name"`
	LocationID     *int64                      `json:"locationId" gorm:"column:location_id"`
	LocationName   *string                     `json:"locationName" gorm:"column:location_name"`
	DurationDays   int                         `json:"durationDays" gorm:"-"`
	DurationNights int                         `json:"durationNights" gorm:"-"`
	Duration       string                      `json:"duration" gorm:"-"`
	SuitableFor    *string                     `json:"suitableFor" gorm:"column:suitable_for"`
	ActivityCount  int                         `json:"activityCount" gorm:"-"`
	Description    *string                     `json:"description" gorm:"column:description"`
	CoverImagePath *string                     `json:"coverImagePath" gorm:"column:cover_image_path"`
	TravelTipsRaw  *string                     `json:"-" gorm:"column:travel_tips"`
	TravelTips     []string                    `json:"travelTips" gorm:"-"`
	Themes         []ItineraryThemeDto         `json:"themes" gorm:"-"`
	FeaturedImages []ItineraryFeaturedImageDto `json:"featuredImages" gorm:"-"`
	Days           []ItineraryDayDto           `json:"days" gorm:"-"`
}

// ItineraryThemeDto represents a theme linked to an itinerary.
type ItineraryThemeDto struct {
	ID   int64  `json:"id" gorm:"column:id"`
	Name string `json:"name" gorm:"column:name"`
}

// ItineraryFeaturedImageDto represents a featured image of an itinerary.
type ItineraryFeaturedImageDto struct {
	ID          int64  `json:"id" gorm:"column:id"`
	ItineraryID int64  `json:"itineraryId" gorm:"column:itinerary_id"`
	ImagePath   string `json:"imagePath" gorm:"column:image_path"`
}

// ItineraryDayDto represents a single day in an itinerary.
type ItineraryDayDto struct {
	ID          int64                  `json:"id" gorm:"column:id"`
	ItineraryID int64                  `json:"itineraryId" gorm:"column:itinerary_id"`
	DayNumber   int                    `json:"dayNumber" gorm:"column:day_number"`
	Title       string                 `json:"title" gorm:"column:title"`
	Subtitle    *string                `json:"subtitle" gorm:"column:subtitle"`
	Summary     *string                `json:"summary" gorm:"column:summary"`
	ImagePath   *string                `json:"imagePath" gorm:"column:image_path"`
	DayTipsRaw  *string                `json:"-" gorm:"column:day_tips"`
	DayTips     []string               `json:"dayTips" gorm:"-"`
	Activities  []ItineraryActivityDto `json:"activities" gorm:"-"`
}

// ItineraryActivityDto represents a single activity in an itinerary day.
type ItineraryActivityDto struct {
	ID               int64   `json:"id" gorm:"column:id"`
	ItineraryDayID   int64   `json:"itineraryDayId" gorm:"column:itinerary_day_id"`
	StartTime        string  `json:"startTime" gorm:"column:start_time"`
	Title            string  `json:"title" gorm:"column:title"`
	ActivityCategory *string `json:"activityCategory" gorm:"column:activity_category"`
	DurationMinutes  *int    `json:"durationMinutes" gorm:"column:duration_minutes"`
	Location         *string `json:"location" gorm:"column:location"`
	GoogleMapsURL    *string `json:"googleMapsUrl" gorm:"column:google_maps_url"`
	Description      *string `json:"description" gorm:"column:description"`
	Notes            *string `json:"notes" gorm:"column:notes"`
	ImagePath        *string `json:"imagePath" gorm:"column:image_path"`
}
