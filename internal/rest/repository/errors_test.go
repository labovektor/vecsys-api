package repository

import (
	"errors"
	"fmt"
	"testing"

	"gorm.io/gorm"
)

func TestMapNotFound(t *testing.T) {
	wrapped := MapNotFound(fmt.Errorf("find participant: %w", gorm.ErrRecordNotFound))
	if !errors.Is(wrapped, ErrNotFound) {
		t.Fatal("wrapped gorm not-found must match ErrNotFound")
	}
	if !errors.Is(wrapped, gorm.ErrRecordNotFound) {
		t.Fatal("original error must stay in the chain for logging")
	}

	if err := MapNotFound(errors.New("boom")); errors.Is(err, ErrNotFound) {
		t.Fatal("unrelated errors must not match ErrNotFound")
	}

	if err := MapNotFound(nil); err != nil {
		t.Fatal("nil must stay nil")
	}
}
