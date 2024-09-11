package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			hash := sha256.Sum256([]byte(os.Getenv("TODO_PASSWORD")))
			hashString := hex.EncodeToString(hash[:])

			var cookieToken string
			cookie, err := r.Cookie("token")
			if err != nil {
				AuthError(w, err)
				return
			} else {
				cookieToken = cookie.Value
				jwtToken, err := jwt.Parse(cookieToken, func(t *jwt.Token) (interface{}, error) {
					return []byte(pass), nil
				})
				if err != nil {
					AuthError(w, err)
					return
				}

				res, ok := jwtToken.Claims.(jwt.MapClaims)
				if !ok {
					AuthError(w, fmt.Errorf("ошибка приведения значени поля Claims к типу wt.MapClaims"))
					return
				}

				hashRaw := res["hash"]

				tokenHash, ok := hashRaw.(string)
				if !ok {
					AuthError(w, fmt.Errorf("ошибка приведения значения хэша к типу string"))
					return
				}

				if tokenHash != hashString {
					AuthError(w, err)
					return
				}

				log.Println("аутентификация пройдена")
			}
		}

		next(w, r)
	}
}

func AuthError(w http.ResponseWriter, err error) {
	log.Printf("ошибка аутентификации: %s", err)
	http.Error(w, "Authentication required", http.StatusUnauthorized)
}
