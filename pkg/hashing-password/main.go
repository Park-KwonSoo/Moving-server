package hashingpassword

import (
	"golang.org/x/crypto/bcrypt"
)

//해쉬패스워드 생성
func GenerateHashPassword(p string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

//해쉬 패스워드와 비교 : hash password & string password
func CompareHashPassword(h, p string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(p))
	if err != nil {
		return false, err
	}
	return true, nil
}
