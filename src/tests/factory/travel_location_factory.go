package factory

import (
	"time"

	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
)

// TravelLocationFactory is a factory for generating TravelLocation dbmodel instances for tests and seeders.
//
// Usage:
//
//	// Create a TravelLocation with default values
//	location := factory.NewTravelLocationFactory().Build()
//
//	// Create a TravelLocation with custom values
//	location := factory.NewTravelLocationFactory().
//	    WithCountryCode("VN").
//	    WithCountryName("Vietnam").
//	    WithLocationName("Đà Lạt").
//	    WithCoordinates(11.9404192, 108.4583132).
//	    WithVisited().
//	    Build()
//
//	// Insert into DB
//	loc, err := factory.NewTravelLocationFactory().Insert(db)
//	// or
//	loc := factory.NewTravelLocationFactory().MustInsert(db)
type TravelLocationFactory struct {
	countryCode   string
	countryName   string
	locationName  string
	latitude      *float64
	longitude     *float64
	imagePath     *string
	status        string
	description   *string
	notes         *string
	googleMapsURL *string
}

// NewTravelLocationFactory creates a new TravelLocationFactory with default values
func NewTravelLocationFactory() *TravelLocationFactory {
	lat := 11.9404192
	lng := 108.4583132
	img := "/images/locations/da_lat.jpg"
	desc := "City of Eternal Spring, famous for cool climate and pine forests."
	notes := "Visited Xuan Huong Lake and local markets."
	mapsURL := "https://maps.google.com/?q=Da+Lat+Lam+Dong+Vietnam"

	return &TravelLocationFactory{
		countryCode:   "VN",
		countryName:   "Vietnam",
		locationName:  "Đà Lạt",
		latitude:      &lat,
		longitude:     &lng,
		imagePath:     &img,
		status:        "visited",
		description:   &desc,
		notes:         &notes,
		googleMapsURL: &mapsURL,
	}
}

// WithLocation sets values from an existing dbmodel.TravelLocation
func (f *TravelLocationFactory) WithLocation(loc dbmodel.TravelLocation) *TravelLocationFactory {
	f.countryCode = loc.CountryCode
	f.countryName = loc.CountryName
	f.locationName = loc.LocationName
	f.latitude = loc.Latitude
	f.longitude = loc.Longitude
	f.imagePath = loc.ImagePath
	f.status = loc.Status
	f.description = loc.Description
	f.notes = loc.Notes
	f.googleMapsURL = loc.GoogleMapsURL
	return f
}

// WithCountryCode sets the ISO 3166-1 alpha-2 country code
func (f *TravelLocationFactory) WithCountryCode(code string) *TravelLocationFactory {
	f.countryCode = code
	return f
}

// WithCountryName sets the country name
func (f *TravelLocationFactory) WithCountryName(name string) *TravelLocationFactory {
	f.countryName = name
	return f
}

// WithLocationName sets the location name
func (f *TravelLocationFactory) WithLocationName(name string) *TravelLocationFactory {
	f.locationName = name
	return f
}

// WithCoordinates sets latitude and longitude
func (f *TravelLocationFactory) WithCoordinates(lat, lng float64) *TravelLocationFactory {
	f.latitude = &lat
	f.longitude = &lng
	return f
}

// WithImagePath sets the image path
func (f *TravelLocationFactory) WithImagePath(path string) *TravelLocationFactory {
	f.imagePath = &path
	return f
}

// WithStatus sets the status ('visited' or 'wishlist')
func (f *TravelLocationFactory) WithStatus(status string) *TravelLocationFactory {
	f.status = status
	return f
}

// WithVisited sets status to 'visited'
func (f *TravelLocationFactory) WithVisited() *TravelLocationFactory {
	f.status = "visited"
	return f
}

// WithWishlist sets status to 'wishlist'
func (f *TravelLocationFactory) WithWishlist() *TravelLocationFactory {
	f.status = "wishlist"
	return f
}

// WithDescription sets the description
func (f *TravelLocationFactory) WithDescription(desc string) *TravelLocationFactory {
	f.description = &desc
	return f
}

// WithNotes sets travel notes
func (f *TravelLocationFactory) WithNotes(notes string) *TravelLocationFactory {
	f.notes = &notes
	return f
}

// WithGoogleMapsURL sets the Google Maps URL
func (f *TravelLocationFactory) WithGoogleMapsURL(url string) *TravelLocationFactory {
	f.googleMapsURL = &url
	return f
}

// Build creates a new dbmodel.TravelLocation instance
func (f *TravelLocationFactory) Build() *dbmodel.TravelLocation {
	now := time.Now()
	return &dbmodel.TravelLocation{
		CountryCode:   f.countryCode,
		CountryName:   f.countryName,
		LocationName:  f.locationName,
		Latitude:      f.latitude,
		Longitude:     f.longitude,
		ImagePath:     f.imagePath,
		Status:        f.status,
		Description:   f.description,
		Notes:         f.notes,
		GoogleMapsURL: f.googleMapsURL,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// BuildWithID creates a new dbmodel.TravelLocation instance with a specific ID
func (f *TravelLocationFactory) BuildWithID(id int64) *dbmodel.TravelLocation {
	loc := f.Build()
	loc.ID = id
	return loc
}

// Insert persists the TravelLocation to the database
func (f *TravelLocationFactory) Insert(db *gorm.DB) (*dbmodel.TravelLocation, error) {
	loc := f.Build()
	if err := db.Create(loc).Error; err != nil {
		return nil, err
	}
	return loc, nil
}

// MustInsert persists the TravelLocation to the database, panicking on error
func (f *TravelLocationFactory) MustInsert(db *gorm.DB) *dbmodel.TravelLocation {
	loc, err := f.Insert(db)
	if err != nil {
		panic(err)
	}
	return loc
}

