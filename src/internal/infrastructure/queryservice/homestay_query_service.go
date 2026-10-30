package queryservice

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/dbutil"
	"platform-go/internal/infrastructure/qsdto"
)

type HomestayQueryService struct {
	db *gorm.DB
}

func NewHomestayQueryService(db *gorm.DB) *HomestayQueryService {
	return &HomestayQueryService{db: db}
}

func (qs *HomestayQueryService) FetchHomestays(ctx context.Context, page, perPage int) (*dbutil.Paginator[qsdto.HomestayDto], error) {
	query := qs.db.WithContext(ctx).Model(&dbmodel.Homestay{}).Order("id ASC")
	paginated, err := dbutil.ModelPaginate[qsdto.HomestayDto](query, perPage, page)
	if err != nil {
		return nil, err
	}
	if len(paginated.Data) == 0 {
		return paginated, nil
	}

	homestayIDs := make([]int64, len(paginated.Data))
	for i, h := range paginated.Data {
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
			imagesByRoomID[img.RoomID] = append(imagesByRoomID[img.RoomID], img)
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

	for i := range paginated.Data {
		if rms, ok := roomsByHomestayID[paginated.Data[i].ID]; ok {
			paginated.Data[i].Rooms = rms
		} else {
			paginated.Data[i].Rooms = []qsdto.HomestayRoomDto{}
		}
	}

	return paginated, nil
}

func (qs *HomestayQueryService) FetchHomestayByID(ctx context.Context, id int64) (*qsdto.HomestayDto, error) {
	var homestay qsdto.HomestayDto
	err := qs.db.WithContext(ctx).Table("homestays").Where("id = ?", id).First(&homestay).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var rooms []qsdto.HomestayRoomDto
	if err := qs.db.WithContext(ctx).Table("homestay_rooms").
		Where("homestay_id = ?", id).
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
			imagesByRoomID[img.RoomID] = append(imagesByRoomID[img.RoomID], img)
		}

		for i := range rooms {
			if imgs, ok := imagesByRoomID[rooms[i].ID]; ok {
				rooms[i].Images = imgs
			} else {
				rooms[i].Images = []qsdto.HomestayImageDto{}
			}
		}
	}

	homestay.Rooms = rooms
	return &homestay, nil
}
