package main

import (
	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(email string, password string) (string, error) {
	combined := email + ":" + password
	bytes, err := bcrypt.GenerateFromPassword([]byte(combined), bcrypt.DefaultCost)
	return string(bytes), err
}

func ValidatePasswordHash(incoming, validator string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(validator), []byte(incoming))
	return err == nil
}
