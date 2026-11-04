package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"platform-go/internal/domain/homestay"
	"platform-go/internal/infrastructure/dbmodel"
)

type homestayRepository struct {
	db DBInterface
}

func NewHomestayRepository(db DBInterface) homestay.HomestayRepository {
	return &homestayRepository{db: db}
}

func (r *homestayRepository) FindHomestayByID(ctx context.Context, id int64) (*homestay.Homestay, error) {
	var model dbmodel.Homestay
	if err := r.db.GetDBByCtx(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &homestay.Homestay{
		ID:            model.ID,
		Name:          model.Name,
		Address:       model.Address,
		Description:   model.Description,
		GoogleMapsURL: model.GoogleMapsURL,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}, nil
}

func (r *homestayRepository) FindRoomsByHomestayID(ctx context.Context, homestayID int64) ([]homestay.HomestayRoom, error) {
	var models []dbmodel.HomestayRoom
	if err := r.db.GetDBByCtx(ctx).Where("homestay_id = ?", homestayID).Find(&models).Error; err != nil {
		return nil, err
	}

	rooms := make([]homestay.HomestayRoom, len(models))
	for i, m := range models {
		rooms[i] = homestay.HomestayRoom{
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

func (r *homestayRepository) FindImagesByRoomID(ctx context.Context, roomID int64) ([]homestay.HomestayImage, error) {
	var models []dbmodel.HomestayImage
	if err := r.db.GetDBByCtx(ctx).Where("room_id = ?", roomID).Find(&models).Error; err != nil {
		return nil, err
	}

	images := make([]homestay.HomestayImage, len(models))
	for i, m := range models {
		images[i] = homestay.HomestayImage{
			ID:        m.ID,
			RoomID:    m.RoomID,
			ImagePath: m.ImagePath,
			SortOrder: m.SortOrder,
			IsCover:   m.IsCover,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}
	return images, nil
}

func (r *homestayRepository) CreateHomestay(ctx context.Context, h *homestay.Homestay) error {
	model := dbmodel.Homestay{
		Name:          h.Name,
		Address:       h.Address,
		Description:   h.Description,
		GoogleMapsURL: h.GoogleMapsURL,
	}
	if err := r.db.GetDBByCtx(ctx).Create(&model).Error; err != nil {
		return err
	}
	h.ID = model.ID
	h.CreatedAt = model.CreatedAt
	h.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *homestayRepository) CreateHomestayRoom(ctx context.Context, room *homestay.HomestayRoom) error {
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

func (r *homestayRepository) CreateHomestayImage(ctx context.Context, images []homestay.HomestayImage) error {
	if len(images) == 0 {
		return nil
	}

	models := make([]dbmodel.HomestayImage, len(images))
	for i, img := range images {
		models[i] = dbmodel.HomestayImage{
			RoomID:    img.RoomID,
			ImagePath: img.ImagePath,
			SortOrder: img.SortOrder,
			IsCover:   img.IsCover,
		}
	}

	if err := r.db.GetDBByCtx(ctx).Create(&models).Error; err != nil {
		return err
	}

	for i := range models {
		images[i].ID = models[i].ID
		images[i].CreatedAt = models[i].CreatedAt
		images[i].UpdatedAt = models[i].UpdatedAt
	}

	return nil
}

func (r *homestayRepository) UpdateHomestay(ctx context.Context, h *homestay.Homestay) error {
	updates := map[string]interface{}{
		"name":            h.Name,
		"address":         h.Address,
		"description":     h.Description,
		"google_maps_url": h.GoogleMapsURL,
	}
	return r.db.GetDBByCtx(ctx).
		Model(&dbmodel.Homestay{}).
		Where("id = ?", h.ID).
		Updates(updates).Error
}

func (r *homestayRepository) UpdateHomestayRoom(ctx context.Context, room *homestay.HomestayRoom) error {
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

func (r *homestayRepository) UpdateHomestayImage(ctx context.Context, image *homestay.HomestayImage) error {
	updates := map[string]interface{}{
		"image_path": image.ImagePath,
		"sort_order": image.SortOrder,
		"is_cover":   image.IsCover,
	}
	return r.db.GetDBByCtx(ctx).
		Model(&dbmodel.HomestayImage{}).
		Where("id = ? AND room_id = ?", image.ID, image.RoomID).
		Updates(updates).Error
}

func (r *homestayRepository) DeleteHomestayByID(ctx context.Context, id int64) error {
	return r.db.GetDBByCtx(ctx).
		Where("id = ?", id).
		Delete(&dbmodel.Homestay{}).Error
}

func (r *homestayRepository) DeleteRoomsByIDs(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.GetDBByCtx(ctx).
		Where("id IN ?", ids).
		Delete(&dbmodel.HomestayRoom{}).Error
}

func (r *homestayRepository) DeleteImagesByRoomIDs(ctx context.Context, roomIDs []int64) error {
	if len(roomIDs) == 0 {
		return nil
	}
	return r.db.GetDBByCtx(ctx).
		Where("room_id IN ?", roomIDs).
		Delete(&dbmodel.HomestayImage{}).Error
}

func (r *homestayRepository) DeleteImagesByIDs(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.GetDBByCtx(ctx).
		Where("id IN ?", ids).
		Delete(&dbmodel.HomestayImage{}).Error
}
