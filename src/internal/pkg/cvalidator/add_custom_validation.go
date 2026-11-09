package cvalidator

import (
	"context"
	"platform-go/internal/infrastructure/repository"

	"github.com/go-playground/validator/v10"
)

type customValidation struct {
	db repository.DBInterface
}

// AddCustomValidation: Add custom validation to the validator
func AddCustomValidation(v *validator.Validate, db repository.DBInterface) {
	cv := customValidation{db}

	_ = v.RegisterValidationCtx("validateTravelLocationID", cv.validateTravelLocationID())
}

const validateTravelLocationIDTranslation = "The specified travel location ID is invalid."

func (v *customValidation) validateTravelLocationID() validator.FuncCtx {
	return func(ctx context.Context, fl validator.FieldLevel) bool {
		travelLocationID := fl.Field().Uint()
		if travelLocationID == 0 {
			return false
		}
		var count int64
		if err := v.db.GetDBByCtx(ctx).Table("travel_locations").
			Where("id = ?", travelLocationID).
			Where("deleted_at IS NULL").
			Count(&count).Error; err != nil {
			return false
		}
		return count > 0
	}
}
