package travelmemoryphoto

import (
	"errors"
	"path"
	"strings"
	"time"
)

var (
	ErrInvalidImagePath   = errors.New("image path must start with travel_memories/")
	ErrDuplicateImagePath = errors.New("duplicate photo image path")
	ErrInvalidSortOrder   = errors.New("sort order must be greater than or equal to zero")
)

type TravelMemoryPhoto struct {
	ID          int64
	MemoryDayID int64
	ImagePath   string
	SortOrder   int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTravelMemoryPhoto(
	memoryDayID int64,
	imagePath string,
	sortOrder int32,
) (TravelMemoryPhoto, error) {
	trimmedPath := strings.TrimSpace(imagePath)
	if !isValidImagePath(trimmedPath) {
		return TravelMemoryPhoto{}, ErrInvalidImagePath
	}

	if sortOrder < 0 {
		return TravelMemoryPhoto{}, ErrInvalidSortOrder
	}

	return TravelMemoryPhoto{
		MemoryDayID: memoryDayID,
		ImagePath:   trimmedPath,
		SortOrder:   sortOrder,
	}, nil
}

func isValidImagePath(value string) bool {
	const prefix = "travel_memories/"

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
