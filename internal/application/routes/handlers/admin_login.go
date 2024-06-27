package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const EnvVarAdminPasswordName = "ADMIN_PASSWORD" // имя переменной окружения в которой хранится пароль админа
const AdminTokenEXP = time.Hour * 3              // время жизни токена админа
const AdminCookieName = "adminAccessToken"       // название куки для хранения доступа админа
const SecretKEY = "supersecretkey"               // ключ для генерации токена

// Claims — структура утверждений, которая включает стандартные утверждения и
// одно пользовательское AdminPassword
type Claims struct {
	jwt.RegisteredClaims
	AdminPassword string
}

type AdminLoginBodyRequest struct {
	Password string `json:"password"`
}

// Обрабатываем запрос на авторизацию админа
func (h *Handler) AdminLogin() http.HandlerFunc {

	// создаем функцию которую будем возвращать как http.HandlerFunc
	fn := func(rw http.ResponseWriter, req *http.Request) {

		var adminLoginBodyRequest AdminLoginBodyRequest

		// после чтения тела запроса, закрываем
		defer req.Body.Close()

		// читаем тело запроса
		body, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		if err = json.Unmarshal(body, &adminLoginBodyRequest); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		// если пароль верный
		adminPassword := os.Getenv(EnvVarAdminPasswordName)
		if adminLoginBodyRequest.Password == adminPassword {

			// генерируем и сохраняем куку
			cookie := createAdminCookie()

			// устанавливаем созданную куку в http
			http.SetCookie(rw, cookie)

			// устанавливаем код 200
			rw.WriteHeader(http.StatusOK)
		} else {
			// устанавливаем код 401
			http.Error(rw, "wrong password", http.StatusUnauthorized)
		}
	}

	return http.HandlerFunc(fn)

}

// создаем куку админа
func createAdminCookie() *http.Cookie {

	// строим строку токена для куки
	token, _ := buildJWTString(os.Getenv(EnvVarAdminPasswordName))

	// создаем куку в http
	cookie := http.Cookie{}
	cookie.Name = AdminCookieName
	cookie.Value = token
	cookie.Expires = time.Now().Add(AdminTokenEXP)
	cookie.Path = "/"

	return &cookie
}

// строим строку для токена
func buildJWTString(pass string) (string, error) {

	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AdminTokenEXP)),
		},
		AdminPassword: pass,
	})

	// получаем ключ для генерации токена
	key := getKeyForTokenGeneration()

	// создаём строку токена
	strToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return strToken, nil
}

// получаем ключ для генерции токена
func getKeyForTokenGeneration() []byte {

	return []byte(SecretKEY)
}
