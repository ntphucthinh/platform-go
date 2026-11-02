package seeder

import (
	"fmt"

	"github.com/hunter-io/geonames"
	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/repository"
	"platform-go/tests/factory"
)

type locationSeedItem struct {
	countryCode   string
	countryName   string
	lookupCity    string
	name          string
	address       string
	latitude      float64
	longitude     float64
	imagePath     string
	status        string
	description   string
	notes         string
	googleMapsURL string
}

// CreateTravelLocations creates initial travel location seed data for 10 Vietnam provinces/cities and several international destinations.
func CreateTravelLocations(db repository.DBInterface) []dbmodel.TravelLocation {
	// Check if travel location data already exists
	var existingLocation dbmodel.TravelLocation
	result := db.GetDB().Where("country_code = ?", "VN").First(&existingLocation)
	if result.Error == nil {
		fmt.Println("Travel location data already exists. Skipping travel location seed.")
		var locations []dbmodel.TravelLocation
		db.GetDB().Find(&locations)
		return locations
	}

	items := []locationSeedItem{
		// 10 Vietnam provinces/cities (includes Đà Lạt, TP. Hồ Chí Minh, Hà Nội, Đà Nẵng, Huế)
		{
			countryCode:   "VN",
			countryName:   "Vietnam",
			lookupCity:    "Da Lat",
			name:          "Đà Lạt",
			address:       "Đà Lạt, Lâm Đồng, Vietnam",
			latitude:      11.9404192,
			longitude:     108.4583132,
			imagePath:     "/images/locations/da_lat.jpg",
			status:        "visited",
			description:   "City of Eternal Spring, famous for cool climate, pine hills, and flower gardens.",
			notes:         "Must visit Xuan Huong Lake, Langbiang Peak, and local night market.",
			googleMapsURL: "https://maps.google.com/?q=Da+Lat+Lam+Dong+Vietnam",
		},
		{
			countryCode:   "VN",
			countryName:   "Vietnam",
			lookupCity:    "Ho Chi Minh City",
			name:          "TP. Hồ Chí Minh",
			address:       "Thành phố Hồ Chí Minh, Vietnam",
			latitude:      10.776889,
			longitude:     106.700806,
			imagePath:     "/images/locations/hcmc.jpg",
			status:        "visited",
			description:   "Dynamic economic hub and largest metropolis of Vietnam.",
			notes:         "Visited Ben Thanh Market, Saigon Opera House, and Nguyen Hue Walking Street.",
			googleMapsURL: "https://maps.google.com/?q=Ho+Chi+Minh+City+Vietnam",
		},
		{
			countryCode:   "VN",
			countryName:   "Vietnam",
			lookupCity:    "Hanoi",
			name:          "Hà Nội",
			address:       "Hoàn Kiếm, Hà Nội, Vietnam",
			latitude:      21.028511,
			longitude:     105.854444,
			imagePath:     "/images/locations/hanoi.jpg",
			status:        "visited",
			description:   "Historic capital city renowned for centuries-old architecture and rich culture.",
			notes:         "Strolling around Hoan Kiem Lake and enjoying egg coffee in the Old Quarter.",
			googleMapsURL: "https://maps.google.com/?q=Hanoi+Vietnam",
		},
		{
			countryCode:   "VN",
			countryName:   "Vietnam",
			lookupCity:    "Da Nang",
			name:          "Đà Nẵng",
			address:       "Đà Nẵng, Vietnam",
			latitude:      16.054407,
			longitude:     108.202167,
			imagePath:     "/images/locations/da_nang.jpg",
			status:        "visited",
			description:   "Coastal city renowned for sandy beaches, Ba Na Hills, and the Dragon Bridge.",
			notes:         "Dragon Bridge fire breathing on weekend nights and My Khe beach swimming.",
			googleMapsURL: "https://maps.google.com/?q=Da+Nang+Vietnam",
		},
		{
			countryCode:   "VN",
			countryName:   "Vietnam",
			lookupCity:    "Hue",
			name:          "Huế",
			address:       "Huế, Thừa Thiên Huế, Vietnam",
			latitude:      16.463713,
			longitude:     107.590866,
			imagePath:     "/images/locations/hue.jpg",
			status:        "visited",
			description:   "Former imperial capital of the Nguyen Dynasty with historic citadels and royal tombs.",
			notes:         "Imperial City tour and boat trip on Perfume River with royal court music.",
			googleMapsURL: "https://maps.google.com/?q=Hue+Vietnam",
		},
		{
			countryCode:   "VN",
			countryName:   "Vietnam",
			lookupCity:    "Nha Trang",
			name:          "Nha Trang",
			address:       "Nha Trang, Khánh Hòa, Vietnam",
			latitude:      12.238791,
			longitude:     109.196749,
			imagePath:     "/images/locations/nha_trang.jpg",
			status:        "visited",
			description:   "Renowned coastal resort city famous for diving sites, offshore islands, and seafood.",
			notes:         "Island hopping in Nha Trang Bay and scuba diving at Hon Mun.",
			googleMapsURL: "https://maps.google.com/?q=Nha+Trang+Khanh+Hoa+Vietnam",
		},
		{
			countryCode:   "VN",
			countryName:   "Vietnam",
			lookupCity:    "Can Tho",
			name:          "Cần Thơ",
			address:       "Ninh Kiều, Cần Thơ, Vietnam",
			latitude:      10.045162,
			longitude:     105.746857,
			imagePath:     "/images/locations/can_tho.jpg",
			status:        "wishlist",
			description:   "The beating heart of the Mekong Delta, famous for Cai Rang floating market and fruit orchards.",
			notes:         "Plan early morning trip to Cai Rang floating market.",
			googleMapsURL: "https://maps.google.com/?q=Can+Tho+Vietnam",
		},
		{
			countryCode:   "VN",
			countryName:   "Vietnam",
			lookupCity:    "Ha Long",
			name:          "Hạ Long",
			address:       "Hạ Long, Quảng Ninh, Vietnam",
			latitude:      20.950454,
			longitude:     107.073364,
			imagePath:     "/images/locations/ha_long.jpg",
			status:        "visited",
			description:   "UNESCO World Heritage site known for emerald waters and thousands of towering limestone islands.",
			notes:         "Overnight cruise on Ha Long Bay with kayaking through Sung Sot cave.",
			googleMapsURL: "https://maps.google.com/?q=Ha+Long+Quang+Ninh+Vietnam",
		},
		{
			countryCode:   "VN",
			countryName:   "Vietnam",
			lookupCity:    "Phu Quoc",
			name:          "Phú Quốc",
			address:       "Phú Quốc, Kiên Giang, Vietnam",
			latitude:      10.289879,
			longitude:     103.98402,
			imagePath:     "/images/locations/phu_quoc.jpg",
			status:        "wishlist",
			description:   "Tropical paradise island known for white-sand beaches, sunset viewpoints, and coral reefs.",
			notes:         "Visit Sao Beach and take the cable car to Hon Thom island.",
			googleMapsURL: "https://maps.google.com/?q=Phu+Quoc+Kien+Giang+Vietnam",
		},
		{
			countryCode:   "VN",
			countryName:   "Vietnam",
			lookupCity:    "Sa Pa",
			name:          "Sa Pa",
			address:       "Sa Pa, Lào Cai, Vietnam",
			latitude:      22.33636,
			longitude:     103.84378,
			imagePath:     "/images/locations/sa_pa.jpg",
			status:        "visited",
			description:   "Mountain town overlooking picturesque terraced rice fields and Mount Fansipan.",
			notes:         "Summit Fansipan peak and trek through Muong Hoa valley.",
			googleMapsURL: "https://maps.google.com/?q=Sa+Pa+Lao+Cai+Vietnam",
		},
		// International destinations
		{
			countryCode:   "JP",
			countryName:   "Japan",
			lookupCity:    "Tokyo",
			name:          "Tokyo",
			address:       "Shibuya, Tokyo, Japan",
			latitude:      35.6762,
			longitude:     139.6503,
			imagePath:     "/images/locations/tokyo.jpg",
			status:        "wishlist",
			description:   "Ultra-modern capital of Japan blending neon-lit skyscrapers with historic temples.",
			notes:         "Visit Shibuya Crossing, Meiji Shrine, and Akihabara.",
			googleMapsURL: "https://maps.google.com/?q=Tokyo+Japan",
		},
		{
			countryCode:   "FR",
			countryName:   "France",
			lookupCity:    "Paris",
			name:          "Paris",
			address:       "Paris, France",
			latitude:      48.8566,
			longitude:     2.3522,
			imagePath:     "/images/locations/paris.jpg",
			status:        "wishlist",
			description:   "Global center for art, fashion, gastronomy, and culture along the River Seine.",
			notes:         "Eiffel Tower summit, Louvre Museum, and Montmartre stroll.",
			googleMapsURL: "https://maps.google.com/?q=Paris+France",
		},
		{
			countryCode:   "US",
			countryName:   "United States",
			lookupCity:    "New York",
			name:          "New York",
			address:       "Manhattan, New York, NY, USA",
			latitude:      40.7128,
			longitude:     -74.0060,
			imagePath:     "/images/locations/new_york.jpg",
			status:        "wishlist",
			description:   "The Big Apple, featuring Times Square, Central Park, and world-class museums.",
			notes:         "Walk the Brooklyn Bridge and watch a Broadway musical.",
			googleMapsURL: "https://maps.google.com/?q=New+York+City+USA",
		},
		{
			countryCode:   "GB",
			countryName:   "United Kingdom",
			lookupCity:    "London",
			name:          "London",
			address:       "Westminster, London, UK",
			latitude:      51.5074,
			longitude:     -0.1278,
			imagePath:     "/images/locations/london.jpg",
			status:        "wishlist",
			description:   "Historic capital on the River Thames, home to iconic landmarks and royal palaces.",
			notes:         "Big Ben, Tower Bridge, and British Museum.",
			googleMapsURL: "https://maps.google.com/?q=London+United+Kingdom",
		},
		{
			countryCode:   "KR",
			countryName:   "South Korea",
			lookupCity:    "Seoul",
			name:          "Seoul",
			address:       "Jongno-gu, Seoul, South Korea",
			latitude:      37.5665,
			longitude:     126.9780,
			imagePath:     "/images/locations/seoul.jpg",
			status:        "wishlist",
			description:   "Vibrant capital known for modern skyscrapers, K-pop culture, and royal palaces.",
			notes:         "Explore Gyeongbokgung Palace and try street food in Myeongdong.",
			googleMapsURL: "https://maps.google.com/?q=Seoul+South+Korea",
		},
	}

	var models []dbmodel.TravelLocation
	for _, item := range items {
		// Resolve state code using geonames package
		stateResult, err := geonames.ResolveState(item.countryCode, item.lookupCity)
		stateCode := ""
		if err == nil && stateResult.State != "" {
			stateCode = stateResult.State
		}

		builder := factory.NewTravelLocationFactory().
			WithCountryCode(item.countryCode).
			WithCountryName(item.countryName).
			WithStateCode(stateCode).
			WithName(item.name).
			WithAddress(item.address).
			WithCoordinates(item.latitude, item.longitude).
			WithImagePath(item.imagePath).
			WithStatus(item.status).
			WithDescription(item.description).
			WithNotes(item.notes).
			WithGoogleMapsURL(item.googleMapsURL)

		models = append(models, *builder.Build())
	}

	// Batch insert into DB
	if len(models) > 0 {
		if err := db.GetDB().Session(&gorm.Session{SkipHooks: true}).Create(&models).Error; err != nil {
			panic(fmt.Sprintf("Failed to create travel location seed data: %v", err))
		}
		fmt.Printf("Seeded %d travel locations successfully.\n", len(models))
	}

	return models
}

