package handlers

import (
	"net/http"
	"time"

	"github.com/cfif1982/taxi/internal/domain/drivers"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const DriverTokenEXP = time.Hour * 3         // время жизни токена водителя
const DriverCookieName = "driverAccessToken" // название куки для хранения доступа водителя
const SecretKEY = "supersecretkey"           // ключ для генерации токена

// Claims — структура утверждений, которая включает стандартные утверждения и
// одно пользовательское DriverID
type Claims struct {
	jwt.RegisteredClaims
	DriverID uuid.UUID
}

type DriverLoginBodyRequest struct {
	Telephone string `json:"telephone"`
	Password  string `json:"password"`
}

// Обрабатываем запрос на авторизацию водителя
func (h *Handler) DriverLogin() http.HandlerFunc {

	// создаем функцию которую будем возвращать как http.HandlerFunc
	fn := func(rw http.ResponseWriter, req *http.Request) {

		var driverRequest DriverLoginBodyRequest

		// считываем запрос в dto
		if err := readRequestToDTO(h, req, &driverRequest); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		// нужно получить пароль водителя по его телефону
		driver, err := h.driverRepo.GetDriverByTelephone(driverRequest.Telephone)
		if err != nil {
			http.Error(rw, drivers.ErrWrongPassword.Error(), http.StatusUnauthorized)
			return
		}

		// Если пароли не совпадают, то алярм
		if driverRequest.Password != driver.Password() {
			http.Error(rw, drivers.ErrWrongPassword.Error(), http.StatusUnauthorized)
			return
		}

		// генерируем куку
		cookie := createDriverCookie(driver.ID())

		// устанавливаем созданную куку в http
		http.SetCookie(rw, cookie)

		// устанавливаем код 200
		rw.WriteHeader(http.StatusOK)
	}

	return http.HandlerFunc(fn)

}

// создаем куку водителя
func createDriverCookie(driverID uuid.UUID) *http.Cookie {

	// строим строку токена для куки
	token, _ := buildJWTString(driverID)

	// создаем куку в http
	cookie := http.Cookie{}
	cookie.Name = DriverCookieName
	cookie.Value = token
	cookie.Expires = time.Now().Add(DriverTokenEXP)
	cookie.Path = "/"

	return &cookie
}

// строим строку для токена
func buildJWTString(driverID uuid.UUID) (string, error) {

	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(DriverTokenEXP)),
		},
		DriverID: driverID,
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
