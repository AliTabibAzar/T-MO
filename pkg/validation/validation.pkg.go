package validation

import (
	"errors"

	"github.com/AliTr404/T-MO/internal/dto"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type (
	Validation struct {
		validator *validator.Validate
	}
)

func NewValidation() echo.Validator {
	return &Validation{validator: validator.New()}
}

func (v *Validation) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "این فیلد الزامی است"
	case "email":
		return "ایمیل غیر معتبر است"
	case "min":
		return "مقدار این فیلد کمتر از مقدار مورد نیاز است"
	}
	return fe.Tag()
}

func ValidateRequest(c echo.Context, i interface{}) []dto.ApiError {

	if err := c.Validate(i); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]dto.ApiError, len(ve))
			for i, fe := range ve {
				out[i] = dto.ApiError{
					Param:   fe.Field(),
					Message: msgForTag(fe),
				}
			}
			return out
		}
	}
	return nil
}
