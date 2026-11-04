package homestay

import (
	"strings"
	"time"
)

type Homestay struct {
	ID            int64
	Name          string
	Address       string
	Description   *string
	GoogleMapsURL *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewHomestay(name, address string, description, googleMapsURL *string) Homestay {
	return Homestay{
		Name:          strings.TrimSpace(name),
		Address:       strings.TrimSpace(address),
		Description:   description,
		GoogleMapsURL: googleMapsURL,
	}
}

func (h *Homestay) Update(name, address string, description, googleMapsURL *string) {
	h.Name = strings.TrimSpace(name)
	h.Address = strings.TrimSpace(address)
	h.Description = description
	h.GoogleMapsURL = googleMapsURL
}

type HomestayRoom struct {
	ID          int64
	HomestayID  int64
	Name        string
	Description *string
	PriceAmount *float64
	PriceUnit   *string
	SortOrder   int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewHomestayRoom(
	homestayID int64,
	name string,
	description *string,
	priceAmount *float64,
	priceUnit *string,
	sortOrder int32,
) HomestayRoom {
	return HomestayRoom{
		HomestayID:  homestayID,
		Name:        strings.TrimSpace(name),
		Description: description,
		PriceAmount: priceAmount,
		PriceUnit:   priceUnit,
		SortOrder:   sortOrder,
	}
}

type HomestayImage struct {
	ID        int64
	RoomID    int64
	ImagePath string
	SortOrder int32
	IsCover   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewHomestayImage(roomID int64, imagePath *string, sortOrder int32, isCover bool) HomestayImage {
	trimmedPath := ""
	if imagePath != nil {
		trimmedPath = strings.TrimSpace(*imagePath)
	}
	return HomestayImage{
		RoomID:    roomID,
		ImagePath: trimmedPath,
		SortOrder: sortOrder,
		IsCover:   isCover,
	}
}
