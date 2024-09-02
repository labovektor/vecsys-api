package util

import "golang.org/x/crypto/bcrypt"

// password user yang diinput dihashing kan
func HashPassword(password string) (string, error) {
	if len(password) < 6 {
		return "", bcrypt.ErrHashTooShort
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// pengecekan password dalam bentuk hashing
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
