package cvalidator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func CustomTagName(v *validator.Validate) {
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		if name == "" {
			name = fld.Tag.Get("label")
		}
		if name == "" {
			name = fld.Name
		}
		return name
	})
}
