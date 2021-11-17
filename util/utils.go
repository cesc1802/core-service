package util

import (
	"fmt"
	"github.com/cesc1802/core-service/i18n"
	"github.com/cesc1802/core-service/model/app_error"
	"github.com/cesc1802/core-service/util/errorcode"
	"github.com/go-playground/validator/v10"
	"strings"
)

func SliceStringToString(sliceVal []string, sep string) string {
	if sep == "" {
		sep = ","
	}
	return strings.Join(sliceVal, sep)
}

func BoolToBoolString(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

func translateToAppVE(i18n *i18n.I18n, lang string,
	valErrors validator.ValidationErrors, errCode string) []app_error.ValidationErrorField {

	res := make([]app_error.ValidationErrorField, len(valErrors))
	for i, valErr := range valErrors {
		res[i] = app_error.ValidationErrorField{
			Field:        valErr.Field(),
			Tag:          valErr.Tag(),
			ErrorMessage: i18n.MustLocalize(lang, fmt.Sprintf("%v.%v", errCode, valErr.Tag()), nil),
		}
	}
	return res
}

func HandleValidationErrors(language string, i18n *i18n.I18n, valErrors validator.ValidationErrors) *app_error.AppError {
	appErr := app_error.ValidationError(
		i18n.MustLocalize(language, errorcode.COM0005, nil),
		"ErrValidation",
		translateToAppVE(i18n, language, valErrors, errorcode.COM0005),
	)
	return appErr
}
