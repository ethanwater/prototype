package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const cost int = 13

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(hash), err
}

func VerfiyHashPassword(hash, password string) bool {
	status := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return status == nil
}
