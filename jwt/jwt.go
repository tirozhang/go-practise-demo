package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func main() {
	var (
		secKey = "123455" // 签发秘钥
	)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "my server",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, 7)),
	})

	tokenString, err := t.SignedString([]byte(secKey))
	fmt.Println(tokenString, err)

	var data jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(tokenString, &data, func(token *jwt.Token) (interface{}, error) {
		return []byte(secKey), nil
	})
	claims := token.Claims.(*jwt.RegisteredClaims)
	fmt.Println(claims.Issuer)

}
