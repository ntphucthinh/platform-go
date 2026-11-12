package homestayroom

import (
	"strings"
	"time"
)

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

