package api

import (
	"errors"
	"fmt"
	"net/http"

	creepsjwt "github.com/Heavenston/creeps_server/creeps_manager/creeps_jwt"
	"github.com/Heavenston/creeps_server/creeps_manager/model"
	"github.com/golang-jwt/jwt/v5"
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

	claims, err := creepsjwt.Decode(auth)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		if errors.Is(err, jwt.ErrTokenExpired) {
			fmt.Fprintf(w, `{"error": "forbidden", "message": "Token expired"}`)
		} else {
			fmt.Fprintf(w, `{"error": "forbidden", "message": "Invalid token"}`)
		}
		return
	}

	userId := claims.UserId
	rs := db.Where("id = ?", userId).First(&user)
	if rs.RowsAffected == 0 {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "forbidden", "message": "Invalid token"}`))
		return
	}
	return
}
