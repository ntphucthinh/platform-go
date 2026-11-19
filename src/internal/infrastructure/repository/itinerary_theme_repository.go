package repository

import (
	"context"

	itinerarytheme "platform-go/internal/domain/itinerary_theme"
	"platform-go/internal/infrastructure/dbmodel"
)

type ItineraryThemeRepository = itinerarytheme.ItineraryThemeRepository

type itineraryThemeRepository struct {
	db DBInterface
}

func NewItineraryThemeRepository(db DBInterface) itinerarytheme.ItineraryThemeRepository {
	return &itineraryThemeRepository{db: db}
}

func (r *itineraryThemeRepository) CountValidThemesByIDs(ctx context.Context, themeIDs []int64) (int64, error) {
	if len(themeIDs) == 0 {
		return 0, nil
	}

	var count int64
	if err := r.db.GetDBByCtx(ctx).
		Model(&dbmodel.TravelTheme{}).
		Where("id IN ? AND deleted_at IS NULL", themeIDs).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *itineraryThemeRepository) CreateItineraryThemes(ctx context.Context, themes []itinerarytheme.ItineraryTheme) error {
	if len(themes) == 0 {
		return nil
	}

	const batchSize = 100
	models := make([]dbmodel.ItineraryTheme, len(themes))
	for i, th := range themes {
		models[i] = dbmodel.ItineraryTheme{
			ItineraryID: th.ItineraryID,
			ThemeID:     th.ThemeID,
		}
	}

	if err := r.db.GetDBByCtx(ctx).CreateInBatches(&models, batchSize).Error; err != nil {
		return err
	}

	for i := range themes {
		themes[i].ID = models[i].ID
		themes[i].CreatedAt = models[i].CreatedAt
	}

	return nil
}

func (r *itineraryThemeRepository) DeleteThemesByItineraryID(ctx context.Context, itineraryID int64) error {
	return r.db.GetDBByCtx(ctx).
		Where("itinerary_id = ?", itineraryID).
		Delete(&dbmodel.ItineraryTheme{}).Error
}
