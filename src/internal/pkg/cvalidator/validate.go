package cvalidator

import (
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
	return cv.v.Struct(obj)
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
