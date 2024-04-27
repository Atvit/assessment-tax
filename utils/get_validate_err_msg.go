package utils

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

const UnknownErrMsg = "unknown error"

const (
	required = "field %s is required"
	oneof    = "the value of %s must be one of %s"
	gt       = "the value of %s must be greater than %s"
	gte      = "the value of %s must be greater than or equal %s"
	lte      = "the value of %s must be less than or equal %s"
	ltefield = "the value of %s value must be lower than or equal value of field %s"
)

type FieldErr struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getErrMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf(required, fe.Field())
	case "oneof":
		return fmt.Sprintf(oneof, fe.Field(), fe.Param())
	case "gt":
		return fmt.Sprintf(gt, fe.Field(), fe.Param())
	case "gte":
		return fmt.Sprintf(gte, fe.Field(), fe.Param())
	case "lte":
		return fmt.Sprintf(lte, fe.Field(), fe.Param())
	case "ltefield":
		return fmt.Sprintf(ltefield, fe.Field(), fe.Param())
	}

	return UnknownErrMsg
}

func GetValidateErrMsg(e error) interface{} {
	var ve validator.ValidationErrors
	if errors.As(e, &ve) {
		out := make([]FieldErr, len(ve))
		for i, fe := range ve {
			out[i] = FieldErr{
				Field:   fe.Field(),
				Message: getErrMsg(fe),
			}
		}

		return out
	}

	return e
}
