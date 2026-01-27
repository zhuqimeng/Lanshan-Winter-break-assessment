package tokens

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func CheckToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("my_secret_key"), nil
	})
	if err != nil || !token.Valid {
		fmt.Println("JWT无效:", err)
		return "", err
	}
	claims, ok := token.Claims.(*TokenClaims)
	if ok {
		fmt.Println("JWT声明:", claims)
		return claims.Username, nil
	}
	return "", nil
}
