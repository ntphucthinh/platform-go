package cvalidator

import (
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var v *validator.Validate

type customValidator struct {
	v *validator.Validate
}

func (cv *customValidator) ValidateStruct(obj any) error {
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		if value.Elem().Kind() != reflect.Struct {
			return cv.ValidateStruct(value.Elem().Interface())
		}
		return cv.v.Struct(obj)
	case reflect.Struct:
		return cv.v.Struct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(binding.SliceValidationError, 0)
		for i := 0; i < count; i++ {
			if err := cv.ValidateStruct(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			}
		}
		if len(validateRet) == 0 {
			return nil
		}
		return validateRet
	default:
		return nil
	}
}

func (cv *customValidator) Engine() any {
	return cv.v
}

func init() {
	NewValidate()
}

// Create a custom validation
// NOTE: A custom validator is required to pass the context
// Reconsider the implementation if this becomes an issue
func NewValidate() {
	v = validator.New(validator.WithRequiredStructEnabled())
	v.SetTagName("binding")

	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, _ = uni.GetTranslator("en")

	_ = en_translations.RegisterDefaultTranslations(v, trans)
	CustomTagName(v)

	// Connect Gin binding to the same validator instance and translations
	binding.Validator = &customValidator{v: v}
}

func Validator() *validator.Validate {
	return v
}
