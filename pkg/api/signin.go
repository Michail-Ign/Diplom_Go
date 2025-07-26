package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

type Password struct {
	Password string `json:"password"`
}

func signinHandler(w http.ResponseWriter, req *http.Request) {

	bodyBytes, err := io.ReadAll(req.Body)

	if err != nil {

		writeJson(w, fmt.Sprintf("ошибка чтения тела (%v)", err))
		return
	}

	var passw Password

	if err := json.Unmarshal(bodyBytes, &passw); err != nil {

		writeJson(w, fmt.Sprintf("ошибка десериализации (%v)", err))
		return
	}

	todoPassword := os.Getenv("TODO_PASSWORD")
	//todo_password = "123" //пока костыль
	if len(todoPassword) == 0 {
		return
	}

	if todoPassword != passw.Password {

		writeJson(w, "Неверный пароль")
		return
	}

	// получаем подписанный токен
	tokenString, err := createJWTToken(passw.Password)
	if err != nil {

		writeJson(w, "Token generation failed")
		return
	}

	//log.Println("token = ", tokenString)
	res := map[string]string{"token": tokenString}
	writeJson(w, res)
}

func createJWTToken(password string) (string, error) {
	claims := jwt.MapClaims{}
	claims["password"] = password // Можно добавить контрольную сумму или любое другое поле

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом
	signingKey := []byte(os.Getenv("SECRET_KEY"))
	if signingKey == nil {
		return "", errors.New("переменная окружения SECRET_KEY не задана")
	}

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
