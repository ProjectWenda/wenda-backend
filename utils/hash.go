package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashToken(token string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing", err)
	}
	return string(bytes)
}

func CheckHash(token string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
	return err == nil
}
