package middlewares

import (
	"context"
	"net/http"

	"github.com/cfif1982/taxi/internal/application"
	"github.com/cfif1982/taxi/internal/application/drivers/handlers"
	"github.com/cfif1982/taxi/internal/domain/drivers"
	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v4"
)

func DriverAuthMiddleware(h http.Handler) http.HandlerFunc {
	adminAuthFn := func(rw http.ResponseWriter, req *http.Request) {

		// получаем токен из куки
		tokenFromCookie, err := req.Cookie(handlers.DriverCookieName)

		// если такой куки нет, то водитель не авторизован
		if err != nil {
			http.Error(rw, drivers.ErrDriverIsNotAuthorized.Error(), http.StatusUnauthorized)
			return
		}

		// получаем user id из токена
		userID, err := getUserIDFromToken(tokenFromCookie.Value)
		if err != nil {
			http.Error(rw, drivers.ErrCookieError.Error(), http.StatusUnauthorized)
			return
		}

		// создаю контекст для сохранения userID
		ctx := context.WithValue(req.Context(), application.KeyDriverID, userID)

		// обрабатываем запрос с контекстом
		h.ServeHTTP(rw, req.WithContext(ctx))
	}

	return http.HandlerFunc(adminAuthFn)
}

// получаем user id из токена
func getUserIDFromToken(tokenString string) (uuid.UUID, error) {
	claims := &handlers.Claims{}

	// получаем ключ для генерации токена
	key := getKeyForTokenGeneration()

	// парсим из строки токена tokenString в структуру claims
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return key, nil
		})

	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, err
	}

	return claims.DriverID, nil
}
