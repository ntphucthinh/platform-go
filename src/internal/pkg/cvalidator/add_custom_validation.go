package cvalidator

import (
	"context"
	"math"
	"reflect"

	"github.com/go-playground/validator/v10"

	"platform-go/internal/infrastructure/repository"
)

type customValidation struct {
	db repository.DBInterface
}

// AddCustomValidation: Add custom validation to the validator
func AddCustomValidation(v *validator.Validate, db repository.DBInterface) {
	cv := customValidation{db}

	_ = v.RegisterValidationCtx("validateTravelLocationID", cv.validateTravelLocationID())
	_ = v.RegisterValidationCtx("validateTravelLocationIDs", cv.validateTravelLocationID())
	_ = v.RegisterValidationCtx("validateTravelThemeID", cv.validateTravelThemeID())
	_ = v.RegisterValidationCtx("validateTravelThemeIDs", cv.validateTravelThemeID())
}

const (
	validateTravelLocationIDTranslation = "The specified travel location ID is invalid."
	validateTravelThemeIDTranslation    = "The specified travel theme ID is invalid."
)

func (v *customValidation) validateTravelLocationID() validator.FuncCtx {
	return func(ctx context.Context, fl validator.FieldLevel) bool {
		return v.validateEntityIDs(ctx, fl, "travel_locations")
	}
}

func (v *customValidation) validateTravelThemeID() validator.FuncCtx {
	return func(ctx context.Context, fl validator.FieldLevel) bool {
		return v.validateEntityIDs(ctx, fl, "travel_themes")
	}
}

func (v *customValidation) validateEntityIDs(ctx context.Context, fl validator.FieldLevel, tableName string) bool {
	ids, ok := extractIDs(fl.Field())
	if !ok {
		return false
	}
	if len(ids) == 0 {
		return true
	}

	uniqueMap := make(map[int64]struct{}, len(ids))
	idList := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			return false
		}
		if _, exists := uniqueMap[id]; !exists {
			uniqueMap[id] = struct{}{}
			idList = append(idList, id)
		}
	}

	var count int64
	if err := v.db.GetDBByCtx(ctx).Table(tableName).
		Where("id IN ?", idList).
		Where("deleted_at IS NULL").
		Count(&count).Error; err != nil {
		return false
	}

	return count == int64(len(idList))
}

func extractIDs(val reflect.Value) ([]int64, bool) {
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		if val.IsNil() {
			return nil, true
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		id := val.Int()
		if id <= 0 {
			return nil, false
		}
		return []int64{id}, true

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u := val.Uint()
		if u > math.MaxInt64 || u == 0 {
			return nil, false
		}
		return []int64{int64(u)}, true

	case reflect.Slice, reflect.Array:
		n := val.Len()
		if n == 0 {
			return []int64{}, true
		}
		var ids []int64
		for i := 0; i < n; i++ {
			elemIDs, ok := extractIDs(val.Index(i))
			if !ok {
				return nil, false
			}
			ids = append(ids, elemIDs...)
		}
		return ids, true

	case reflect.Struct:
		f := val.FieldByName("ID")
		if !f.IsValid() {
			f = val.FieldByName("Id")
		}
		if f.IsValid() {
			return extractIDs(f)
		}
		return nil, false

	default:
		return nil, false
	}
}
