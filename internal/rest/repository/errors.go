package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// ErrNotFound is the repository contract for a missing record.
// Callers must use errors.Is instead of depending on gorm.ErrRecordNotFound.
var ErrNotFound = errors.New("not found")

// MapNotFound converts gorm's not-found error into ErrNotFound while keeping
// the original error in the chain for logging.
func MapNotFound(err error) error {
	if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return fmt.Errorf("%w: %w", ErrNotFound, err)
}
