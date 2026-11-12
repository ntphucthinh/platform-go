package cvalidator

import (
	"context"
	"math"
	"platform-go/internal/infrastructure/repository"
	"reflect"

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
		var travelLocationID int64
		switch fl.Field().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			travelLocationID = fl.Field().Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			u := fl.Field().Uint()
			if u > math.MaxInt64 {
				return false
			}
			travelLocationID = int64(u)
		default:
			return false
		}
		if travelLocationID <= 0 {
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
