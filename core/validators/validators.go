package validators

import (
	"github.com/go-playground/validator/v10"
	"time"
)

// ValidateAfterDate 判断时间是否是当前时间之后
func ValidateAfterDate(fl validator.FieldLevel) bool {
	date, err := time.Parse("2006-01-02", fl.Field().String())
	if err != nil {
		return false
	}
	if date.Before(time.Now()) {
		return false
	}
	return true
}

func ValidateMobile(fl validator.FieldLevel) bool {
	return true
}

func ValidateEmail(fl validator.FieldLevel) bool {
	return true
}

func ValidateDate(fl validator.FieldLevel) bool {
	return true
}
