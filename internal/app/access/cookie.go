package access

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/KraDM09/shortener/internal/app/util"
	"github.com/KraDM09/shortener/internal/constants"
	"github.com/golang-jwt/jwt/v4"
)

type Cookie struct{}

type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

const SecretKey = "secret"

func (c Cookie) Request(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var userID string
		token, err := r.Cookie("token")
		isNewToken := false

		switch {
		case errors.Is(err, http.ErrNoCookie):
			userID = util.CreateUUID()
			token, err := GenerateJWT(userID)
			isNewToken = true

			if err != nil {
				panic(fmt.Errorf("ошибка при генерации jwt для пользователя без токена %w", err))
			}

			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   token,
				Expires: time.Now().Add(24 * 7 * time.Hour),
			})
		case err != nil:
			panic(fmt.Errorf("ошибка при получении токена из куки %w", err))
		}

		if isNewToken != true {
			userID = GetUserID(token.Value)
		}

		if userID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), constants.ContextUserIDKey, userID)

		h.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func GenerateJWT(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(SecretKey))
}

func GetUserID(tokenString string) string {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SecretKey), nil
		})
	if err != nil {
		return ""
	}

	if !token.Valid {
		fmt.Println("Token is not valid")
		return ""
	}

	fmt.Println("Token os valid")
	return claims.UserID
}
