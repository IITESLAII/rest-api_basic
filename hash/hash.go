package hash

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func PasswordHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error while hashing password: %v", err)
	}
	return string(hash)
}
func CheckPasswordHash(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		log.Println("Error while comparing hash and password: %v", err)
		return false
	}
	return true
}
