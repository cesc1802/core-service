package validator

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

func JsonTagNameFunc(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}

func InitValidator() *validator.Validate {
	validate := validator.New()

	// register function to get tag name from json tags.
	validate.RegisterTagNameFunc(JsonTagNameFunc)

	return validate
}
