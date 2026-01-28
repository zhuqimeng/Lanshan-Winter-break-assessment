package tokens

import (
	"fmt"
	"zhihu/app/api/configs"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func CheckToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("my_secret_key"), nil
	})
	if err != nil || !token.Valid {
		configs.Logger.Warn("CheckToken", zap.Error(err))
		return "", err
	}
	claims, ok := token.Claims.(*TokenClaims)
	if ok {
		configs.Logger.Info("CheckToken", zap.Any("JWT声明:", claims))
		return claims.Username, nil
	}
	return "", nil
}
