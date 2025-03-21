package util

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Mapping pesan error
var validationMessages = map[string]string{
	"required": "harus diisi",
	"email":    "Format email tidak valid",
	"min":      "minimal %s karakter",
	"max":      "maksimal %s karakter",
	"phone":    "harus dimulai dengan 628 dan terdiri dari 10-13 angka",
}

// Custom validator untuk nomor telepon Indonesia (format 628xxx)
func IndonesianPhoneValidator(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	matched, _ := regexp.MatchString(`^628\d{8,11}$`, phone) // 628 + 8-11 digit angka
	return matched
}

func ValidateStruct(s any) error {
	err := validate.Struct(s)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors) // Cast ke ValidationErrors
		firstErr := validationErrors[0]                      // Ambil error pertama

		msg := validationMessages[firstErr.Tag()]
		if msg == "" {
			msg = "input tidak valid"
		}

		if firstErr.Param() != "" {
			msg = fmt.Sprintf(msg, firstErr.Param())
		}

		return fmt.Errorf("%s %s", firstErr.Field(), msg)
	}

	return nil
}

func InitValidator() {
	validate = validator.New()
	validate.RegisterValidation("phone", IndonesianPhoneValidator)
}
