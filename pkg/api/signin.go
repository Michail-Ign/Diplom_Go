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

	var err_ret ErrorRet

	bodyBytes, err := io.ReadAll(req.Body)

	if err != nil {
		//http.StatusInternalServerError
		err_ret.Error = fmt.Sprintf("ошибка чтения тела (%v)", err)
		err = writeJson(w, err_ret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	var passw Password

	if err := json.Unmarshal(bodyBytes, &passw); err != nil {

		err_ret.Error = fmt.Sprintf("ошибка десериализации (%v)", err)
		err = writeJson(w, err_ret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	todo_password := os.Getenv("TODO_PASSWORD")
	//todo_password = "123" //пока костыль
	if len(todo_password) == 0 {
		return
	}

	if todo_password != passw.Password {

		err_ret.Error = "Неверный пароль"
		err := writeJson(w, err_ret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	/*token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	})*/

	// получаем подписанный токен
	tokenString, err := createJWTToken(passw.Password)
	if err != nil {
		//fmt.Printf("failed to sign jwt: %s\n", err)
		err_ret.Error = "Token generation failed"
		err := writeJson(w, err_ret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	//log.Println("token = ", tokenString)

	res := map[string]string{"token": tokenString}
	err = writeJson(w, res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
