package middlewares

import (
	"net/http"

	"github.com/cfif1982/taxi/internal/application/routes/handlers"

	"github.com/golang-jwt/jwt/v4"
)

func AdminAuthMiddleware(h http.Handler) http.HandlerFunc {
	adminAuthFn := func(rw http.ResponseWriter, req *http.Request) {

		// получаем токен из куки
		tokenFromCookie, err := req.Cookie(handlers.AdminCookieName)

		// если такой куки нет, то админ не авторизован
		if err != nil {
			http.Error(rw, "admin not authorized", http.StatusUnauthorized)
			return
		}

		// получаем пароль админа из токена
		adminPas := getAdminPasswordFromToken(tokenFromCookie.Value)

		// если пароль неверный, то админ не авторизован
		if adminPas != handlers.AdminPassword {
			http.Error(rw, "wrong password", http.StatusUnauthorized)
			return
		}

		// обрабатываем сам запрос
		h.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(adminAuthFn)
}

// получаем пароль админа из токена
func getAdminPasswordFromToken(tokenString string) string {
	claims := &handlers.Claims{}

	// получаем ключ для генерации токена
	key := getKeyForTokenGeneration()

	// парсим из строки токена tokenString в структуру claims
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return key, nil
		})

	if err != nil {
		return ""
	}

	if !token.Valid {
		return ""
	}

	return claims.AdminPassword
}

// получаем ключ для генерции токена
func getKeyForTokenGeneration() []byte {

	return []byte(handlers.SecretKEY)
}
