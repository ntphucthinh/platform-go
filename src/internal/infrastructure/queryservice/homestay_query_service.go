package queryservice

import (
	"context"

	"gorm.io/gorm"

	"platform-go/internal/infrastructure/qsdto"
)

type HomestayQueryService struct {
	db *gorm.DB
}

func NewHomestayQueryService(db *gorm.DB) *HomestayQueryService {
	return &HomestayQueryService{db: db}
}

func (qs *HomestayQueryService) FetchHomestays(ctx context.Context) ([]qsdto.HomestayDto, error) {
	var homestays []qsdto.HomestayDto
	if err := qs.db.WithContext(ctx).Table("homestays").Order("id ASC").Find(&homestays).Error; err != nil {
		return nil, err
	}
	if len(homestays) == 0 {
		return []qsdto.HomestayDto{}, nil
	}

	homestayIDs := make([]int64, len(homestays))
	for i, h := range homestays {
		homestayIDs[i] = h.ID
	}

	var rooms []qsdto.HomestayRoomDto
	if err := qs.db.WithContext(ctx).Table("homestay_rooms").
		Where("homestay_id IN ?", homestayIDs).
		Order("sort_order ASC, id ASC").
		Find(&rooms).Error; err != nil {
		return nil, err
	}

	if len(rooms) > 0 {
		roomIDs := make([]int64, len(rooms))
		for i, r := range rooms {
			roomIDs[i] = r.ID
		}

		var images []qsdto.HomestayImageDto
		if err := qs.db.WithContext(ctx).Table("homestay_images").
			Where("room_id IN ?", roomIDs).
			Order("sort_order ASC, id ASC").
			Find(&images).Error; err != nil {
			return nil, err
		}

		imagesByRoomID := make(map[int64][]qsdto.HomestayImageDto)
		for _, img := range images {
			if img.RoomID != nil {
				imagesByRoomID[*img.RoomID] = append(imagesByRoomID[*img.RoomID], img)
			}
		}

		for i := range rooms {
			if imgs, ok := imagesByRoomID[rooms[i].ID]; ok {
				rooms[i].Images = imgs
			} else {
				rooms[i].Images = []qsdto.HomestayImageDto{}
			}
		}
	}

	roomsByHomestayID := make(map[int64][]qsdto.HomestayRoomDto)
	for _, room := range rooms {
		roomsByHomestayID[room.HomestayID] = append(roomsByHomestayID[room.HomestayID], room)
	}

	for i := range homestays {
		if rms, ok := roomsByHomestayID[homestays[i].ID]; ok {
			homestays[i].Rooms = rms
		} else {
			homestays[i].Rooms = []qsdto.HomestayRoomDto{}
		}
	}

	return homestays, nil
}
