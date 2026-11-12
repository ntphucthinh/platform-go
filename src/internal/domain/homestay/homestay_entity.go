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
