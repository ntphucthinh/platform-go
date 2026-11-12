package homestayimage

import (
	"strings"
	"time"
)

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

