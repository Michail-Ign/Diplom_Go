package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Password             string // идентификатор пользователя
	jwt.RegisteredClaims        // базовый тип
}

func Init() {
	http.HandleFunc("/api/nextdate", auth(nextDayHandler))
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/task/done", auth(taskDoneHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/signin", auth(signinHandler))
}

func auth(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")

		if len(pass) > 0 {
			var jwt string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}

			//var valid bool
			// здесь код для валидации и проверки JWT-токена

			var claim Claims
			if !ValidateToken(jwt, &claim) {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

		}
		next(w, r)
	})
}

func ValidateToken(tokenString string, claims *Claims) bool {

	secretKey := os.Getenv("JWT_SECRET")

	// парсим токен
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// проверяем алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неверный метод подписи")
		}
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return false
	}
	return true
}
