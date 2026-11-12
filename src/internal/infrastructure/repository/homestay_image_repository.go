package repository

import (
	"context"

	homestayimage "platform-go/internal/domain/homestay_image"
	"platform-go/internal/infrastructure/dbmodel"
)

type HomestayImageRepository = homestayimage.HomestayImageRepository

type homestayImageRepository struct {
	db DBInterface
}

func NewHomestayImageRepository(db DBInterface) homestayimage.HomestayImageRepository {
	return &homestayImageRepository{db: db}
}

func (r *homestayImageRepository) FindImagesByRoomID(ctx context.Context, roomID int64) ([]homestayimage.HomestayImage, error) {
	var models []dbmodel.HomestayImage
	if err := r.db.GetDBByCtx(ctx).Where("room_id = ?", roomID).Find(&models).Error; err != nil {
		return nil, err
	}

	images := make([]homestayimage.HomestayImage, len(models))
	for i, m := range models {
		images[i] = homestayimage.HomestayImage{
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

func (r *homestayImageRepository) CreateHomestayImages(ctx context.Context, images []homestayimage.HomestayImage) error {
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

func (r *homestayImageRepository) UpdateHomestayImage(ctx context.Context, image *homestayimage.HomestayImage) error {
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

func (r *homestayImageRepository) DeleteImagesByRoomIDs(ctx context.Context, roomIDs []int64) error {
	if len(roomIDs) == 0 {
		return nil
	}
	return r.db.GetDBByCtx(ctx).
		Where("room_id IN ?", roomIDs).
		Delete(&dbmodel.HomestayImage{}).Error
}

func (r *homestayImageRepository) DeleteImagesByIDs(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.GetDBByCtx(ctx).
		Where("id IN ?", ids).
		Delete(&dbmodel.HomestayImage{}).Error
}

