package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"platform-go/internal/domain/homestay"
	"platform-go/internal/infrastructure/dbmodel"
)

type HomestayRepository = homestay.HomestayRepository

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

func (r *homestayRepository) DeleteHomestayByID(ctx context.Context, id int64) error {
	return r.db.GetDBByCtx(ctx).
		Where("id = ?", id).
		Delete(&dbmodel.Homestay{}).Error
}
