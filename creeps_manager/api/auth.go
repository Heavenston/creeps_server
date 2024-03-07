package api

import (
	"fmt"
	"net/http"

	"github.com/Heavenston/creeps_server/creeps_manager/keys"
	"github.com/Heavenston/creeps_server/creeps_manager/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func auth(db *gorm.DB, w http.ResponseWriter, req *http.Request) (user model.User, err error) {
	auth := req.Header.Get("Authorization")
	if auth == "" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "forbidden", "message": "Missing auth header"}`))
		err = fmt.Errorf("No auth header")
		return
	}

	token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return keys.JWTSecret, nil
	})
	if err != nil {
		log.Debug().Err(err).Msg("token parse error")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "forbidden", "message": "Invalid auth header"}`))
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "forbidden", "message": "Invalid auth header"}`))
		return
	}

	userId := claims["uid"]
	rs := db.Where("id = ?", userId).First(&user)
	if rs.RowsAffected == 0 {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "forbidden", "message": "Could not find user"}`))
		return
	}
	return
}
