package util

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
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

// IndonesianPhoneValidator Custom validator untuk nomor telepon Indonesia (format 628xxx)
func IndonesianPhoneValidator(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	matched, _ := regexp.MatchString(`^628\d{8,11}$`, phone) // 628 + 8-11 digit angka
	return matched
}

func ValidateStruct(s any) error {
	err := validate.Struct(s)
	if err != nil {
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors) // Cast ke ValidationErrors
		firstErr := validationErrors[0]   // Ambil error pertama

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

type FileValidationOpts struct {
	MinSize           int64 // dalam byte
	MaxSize           int64 // dalam byte
	AllowedExtensions []string
}

func DefaultFileValidationOpts() *FileValidationOpts {
	return &FileValidationOpts{
		MinSize:           0,                                 // 0MB
		MaxSize:           toBytes(5),                        // 5MB
		AllowedExtensions: []string{".jpg", ".jpeg", ".png"}, // Image
	}
}

func ValidateFile(file *multipart.FileHeader, opts ...*FileValidationOpts) error {
	if file == nil {
		return errors.New("file is required")
	}

	opt := DefaultFileValidationOpts()

	if len(opts) > 0 {
		iOpts := opts[0]
		if iOpts.MinSize > 0 {
			opt.MinSize = iOpts.MinSize
		}

		if iOpts.MaxSize > 0 {
			opt.MaxSize = iOpts.MaxSize
		}

		if iOpts.AllowedExtensions != nil {
			opt.AllowedExtensions = iOpts.AllowedExtensions
		}
	}

	// Validasi ukuran file
	if file.Size < opt.MinSize {
		return errors.New("file size is too small")
	}
	if file.Size > opt.MaxSize {
		return fmt.Errorf("minimum file size is %.1fMB", toMB(opt.MaxSize))
	}

	// Validasi ekstensi
	ext := filepath.Ext(file.Filename)
	allowed := false
	for _, allowedExt := range opt.AllowedExtensions {
		if ext == allowedExt {
			allowed = true
			break
		}
	}

	if !allowed {
		return errors.New("file extension is not allowed")
	}

	return nil
}

func InitValidator() {
	validate = validator.New()
	validate.RegisterValidation("phone", IndonesianPhoneValidator)
}

func toMB(byte int64) float64 {
	return float64(byte) / (1024 * 1024)
}

func toBytes(mb float64) int64 {
	return int64(mb * 1024 * 1024)
}
