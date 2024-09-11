package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"

	"todo/model"
	"todo/service"
)

func (h *TodoHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var login model.Login
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		service.ErrorResponse(w, "body getting error", err)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &login); err != nil {
		service.ErrorResponse(w, "JSON encoding error", err)
		return
	}

	pass := os.Getenv("TODO_PASSWORD")

	if login.Password == pass {
		hash := sha256.Sum256([]byte(login.Password))
		hashString := hex.EncodeToString(hash[:])
		hashClaims := jwt.MapClaims{"hash": hashString}

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, hashClaims)
		token, err := jwtToken.SignedString([]byte(pass))
		if err != nil {
			log.Println(err.Error())
			service.ErrorResponse(w, "внутренняя ошибка сервера", err)
		}

		authData, err := json.Marshal(model.Auth{Token: token})
		if err != nil {
			log.Println(err.Error())
			service.ErrorResponse(w, "ошибка сериализации данных", err)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(authData)
		if err != nil {
			log.Println(err.Error())
			service.ErrorResponse(w, "внутренняя ошибка сервера", err)
			return
		}
	} else {
		errorResponse := model.Error{Message: "неправильный пароль"}
		errorData, _ := json.Marshal(errorResponse)
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write(errorData)
		if err != nil {
			log.Println(err.Error())
			service.ErrorResponse(w, "внутренняя ошибка сервера", err)
			return
		}
	}
}
