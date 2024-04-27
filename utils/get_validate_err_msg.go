package utils

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

const UnknownErrMsg = "unknown error"

type FieldErr struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getErrMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("field %s is required", fe.Field())
	case "oneof":
		return fmt.Sprintf("the value of %s must be one of %s", fe.Field(), fe.Param())
	case "gt":
		return fmt.Sprintf("the value of %s must be greater than %s", fe.Field(), fe.Param())
	case "gte":
		return fmt.Sprintf("the value of %s must be greater than or equal %s", fe.Field(), fe.Param())
	case "ltefield":
		return fmt.Sprintf("the value of %s value must be lower than or equal value of field %s", fe.Field(), fe.Param())
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
