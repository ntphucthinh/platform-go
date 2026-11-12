package repository

import (
	"context"

	homestayroom "platform-go/internal/domain/homestay_room"
	"platform-go/internal/infrastructure/dbmodel"
)

type HomestayRoomRepository = homestayroom.HomestayRoomRepository

type homestayRoomRepository struct {
	db DBInterface
}

func NewHomestayRoomRepository(db DBInterface) homestayroom.HomestayRoomRepository {
	return &homestayRoomRepository{db: db}
}

func (r *homestayRoomRepository) FindRoomsByHomestayID(ctx context.Context, homestayID int64) ([]homestayroom.HomestayRoom, error) {
	var models []dbmodel.HomestayRoom
	if err := r.db.GetDBByCtx(ctx).Where("homestay_id = ?", homestayID).Find(&models).Error; err != nil {
		return nil, err
	}

	rooms := make([]homestayroom.HomestayRoom, len(models))
	for i, m := range models {
		rooms[i] = homestayroom.HomestayRoom{
			ID:          m.ID,
			HomestayID:  m.HomestayID,
			Name:        m.Name,
			Description: m.Description,
			PriceAmount: m.PriceAmount,
			PriceUnit:   m.PriceUnit,
			SortOrder:   m.SortOrder,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
		}
	}
	return rooms, nil
}

func (r *homestayRoomRepository) CreateHomestayRoom(ctx context.Context, room *homestayroom.HomestayRoom) error {
	model := dbmodel.HomestayRoom{
		HomestayID:  room.HomestayID,
		Name:        room.Name,
		Description: room.Description,
		PriceAmount: room.PriceAmount,
		PriceUnit:   room.PriceUnit,
		SortOrder:   room.SortOrder,
	}
	if err := r.db.GetDBByCtx(ctx).Create(&model).Error; err != nil {
		return err
	}
	room.ID = model.ID
	room.CreatedAt = model.CreatedAt
	room.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *homestayRoomRepository) UpdateHomestayRoom(ctx context.Context, room *homestayroom.HomestayRoom) error {
	updates := map[string]interface{}{
		"name":         room.Name,
		"description":  room.Description,
		"price_amount": room.PriceAmount,
		"price_unit":   room.PriceUnit,
		"sort_order":   room.SortOrder,
	}
	return r.db.GetDBByCtx(ctx).
		Model(&dbmodel.HomestayRoom{}).
		Where("id = ? AND homestay_id = ?", room.ID, room.HomestayID).
		Updates(updates).Error
}

func (r *homestayRoomRepository) DeleteRoomsByIDs(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.GetDBByCtx(ctx).
		Where("id IN ?", ids).
		Delete(&dbmodel.HomestayRoom{}).Error
}

