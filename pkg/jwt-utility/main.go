package jwtutility

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtSecretKey string = os.Getenv("JWT_SECRET_KEY")

type accessToken struct {
	jwt.StandardClaims
	UserId string `json:"userId"`
}

//memId를 바탕으로 token을 발급
func GetJwtToken(memId string) (string, error) {
	claims := &accessToken{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(), //30일 토큰
		},
		UserId: memId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecretKey))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

//token의 유효성 검증 및 decode
func ValidateToken(tokenString string) (string, error) {
	claims := &accessToken{}
	rslt, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})

	if !rslt.Valid {
		return "", err
	}

	return claims.UserId, nil
}
