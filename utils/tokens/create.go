package tokens

import (
	"time"
	"zhihu/app/api/configs"

	"github.com/bwmarrin/snowflake"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func MakeToken(username string, expTime time.Time) (string, error) {
	snowflakeNode, err := snowflake.NewNode(1)
	if err != nil {
		return "", err
	}
	tokenID := snowflakeNode.Generate()
	id := tokenID.String()
	claims := TokenClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ZhiHu",
			Subject:   id,
			Audience:  []string{"ZhiHu"},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expTime),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        id,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("my_secret_key"))
	if err != nil {
		configs.Logger.Error("MakeToken", zap.Error(err))
		return "", err
	}
	return tokenString, nil
}
