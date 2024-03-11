package creepsjwt

import (
	"fmt"
	"time"

	"github.com/Heavenston/creeps_server/creeps_manager/keys"
	"github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	UserId int `json:"uid"`
}

type Claims struct {
	jwt.RegisteredClaims
	Payload
}

func Encode(userId int) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 2)),
		},
		Payload: Payload{
			UserId: userId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	strToken, err := token.SignedString(keys.JWTSecret)
	if err != nil {
		return strToken, err
	}
	return strToken, err
}

func Decode(strToken string) (Claims, error) {
	token, err := jwt.ParseWithClaims(strToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return keys.JWTSecret, nil
	}, jwt.WithIssuedAt(), jwt.WithExpirationRequired())
	if err != nil {
		return Claims{}, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return Claims{}, fmt.Errorf("Invalid claims")
	}

	return *claims, nil
}
