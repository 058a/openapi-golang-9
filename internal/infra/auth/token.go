package auth

import (
	"errors"
	"openapi/internal/infra/env"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtCustomClaims struct {
	UserId uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func EncodeToken(userId uuid.UUID) (string, error) {
	claims := &jwtCustomClaims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 60)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(env.GetJwtSecret()))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func DecodeToken(authorization string) (*jwtCustomClaims, error) {
	unsignedToken := strings.TrimPrefix(authorization, "Bearer ")
	token, err := jwt.Parse(unsignedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, errors.Join(err, err)
	}

	if !token.Valid {
		return nil, errors.Join(err, err)
	}

	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok {
		return nil, errors.New("couldn't parse claims")
	}
	return claims, nil
}
