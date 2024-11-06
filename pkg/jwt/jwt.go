package jwt

import (
	"awesomeProject/pkg/config"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type claims struct {
	Username string `json:"username"`
	uuid     string
	jwt.RegisteredClaims
}

var jwtKey = []byte(config.GetJWTConfig().SecretKey)

func GenerateJWT(username, uuid string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	rc := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}
	claims := &claims{
		Username:         username,
		RegisteredClaims: rc,
		uuid:             uuid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseJWT(tokenString string) (string, error) {
	withClaims, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		fmt.Errorf("Error while parsing JWT: %v", err)
	}

	if claims, ok := withClaims.Claims.(*claims); ok && withClaims.Valid {
		return claims.Username, nil
	} else {
		return "", fmt.Errorf("Invalid JWT")
	}

}
