package middlewares

// import (
// 	"context"
// 	"crypto/rand"
// 	"encoding/binary"
// 	"net/http"
// 	"time"

// 	// "github.com/cfif1982/taxi/internal/application/handlers"

// 	"github.com/golang-jwt/jwt/v4"
// )

// // Claims — структура утверждений, которая включает стандартные утверждения и
// // одно пользовательское UserID
// type Claims struct {
// 	jwt.RegisteredClaims
// 	UserID int
// }

// // создаю свой тип для ключей контекста
// type ctxKey string

// const KeyUserID ctxKey = "user_id" //  ключ в контексте для поля user id

// const TokenEXP = time.Hour * 3        // время жизни токена
// const CookieName = "adminAccessToken" // название куки для хранения доступа админа

// func AuthMiddleware(h http.Handler) http.Handler {
// 	authFn := func(rw http.ResponseWriter, req *http.Request) {

// 		// получаем токен из куки
// 		tokenFromCookie, err := req.Cookie(AdminCookieName)

// 		// если такой куки нет, то создаем куку
// 		if err != nil {
// 			// генерируем user id
// 			userID, err := createUserID()
// 			if err != nil {
// 				http.Error(rw, "user create error", http.StatusInternalServerError)
// 				return
// 			}

// 			cookie := createCookie(userID)

// 			// устанавливаем созданную куку в http
// 			http.SetCookie(rw, cookie)

// 			// тут по заданию сказано, что если при обращении по /api/user/urls нужно проверять условие:
// 			// если не авторизован, то нужно вывести сообщение об этом
// 			// но другое условие - при обращении на добавление: если не авторизован, то создать user id и авторизовать
// 			// я сделал так: при обращении по любому адресу если не авторизован, то создаю user id и авторизую
// 			// но из-за этого приходится тут проверять адрес и делать исключение
// 			// хотя по логике нужно делать так, как ты говорил - отдельным запросом авторизация и уже дальше вся логика
// 			// а так приходится костыль такой делать((
// 			// узнаю адрес запроса
// 			// uri := req.RequestURI

// 			// // если запрос на вывод всех url пользователя, то не сохраняем user id
// 			// if uri != "/api/user/urls" {
// 			// 	// создаю контекст для сохранения userID
// 			// 	ctx := context.WithValue(req.Context(), handlers.KeyUserID, userID)

// 			// 	// обрабатываем запрос
// 			// 	h.ServeHTTP(rw, req.WithContext(ctx))
// 			// } else {
// 			// 	// обрабатываем запрос
// 			// 	h.ServeHTTP(rw, req)
// 			// }

// 			return
// 		}

// 		// получаем user id из токена
// 		userID := getUserIDFromToken(tokenFromCookie.Value)

// 		// если в токенен нет узера, то заново создаем куку
// 		if userID == -1 {
// 			// генерируем user id
// 			userID, err := createUserID()
// 			if err != nil {
// 				http.Error(rw, "user create error", http.StatusInternalServerError)
// 				return
// 			}

// 			cookie := createCookie(userID)

// 			// устанавливаем созданную куку в http
// 			http.SetCookie(rw, cookie)

// 			// создаю контекст для сохранения userID
// 			ctx := context.WithValue(req.Context(), KeyUserID, userID)

// 			// обрабатываем запрос
// 			h.ServeHTTP(rw, req.WithContext(ctx))
// 		}

// 		// создаю контекст для сохранения userID
// 		ctx := context.WithValue(req.Context(), KeyUserID, userID)

// 		// обрабатываем сам запрос
// 		h.ServeHTTP(rw, req.WithContext(ctx))
// 	}

// 	return http.HandlerFunc(authFn)
// }

// // создаем куку
// func createCookie(userID int) *http.Cookie {

// 	// строим строку токена для куки
// 	token, _ := buildJWTString(userID)

// 	// создаем куку в http
// 	cookie := http.Cookie{}
// 	cookie.Name = AdminCookieName
// 	cookie.Value = token
// 	cookie.Expires = time.Now().Add(AdminTokenEXP)
// 	cookie.Path = "/"

// 	return &cookie
// }

// // получаем user id из токена
// func getUserIDFromToken(tokenString string) int {
// 	claims := &Claims{}

// 	// получаем ключ для генерации токена
// 	key := getKeyForTokenGeneration()

// 	// парсим из строки токена tokenString в структуру claims
// 	token, err := jwt.ParseWithClaims(tokenString, claims,
// 		func(t *jwt.Token) (interface{}, error) {
// 			return key, nil
// 		})

// 	if err != nil {
// 		return -1
// 	}

// 	if !token.Valid {
// 		return -1
// 	}

// 	return claims.UserID
// }

// // строим строку для токена
// func buildJWTString(userID int) (string, error) {

// 	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AdminTokenEXP)),
// 		},
// 		UserID: userID,
// 	})

// 	// получаем ключ для генерации токена
// 	key := getKeyForTokenGeneration()

// 	// создаём строку токена
// 	strToken, err := token.SignedString(key)
// 	if err != nil {
// 		return "", err
// 	}

// 	// возвращаем строку токена
// 	return strToken, nil
// }

// // генерируем user id
// func createUserID() (int, error) {

// 	// генерируем случайную последовательность из 6 байт
// 	b := make([]byte, 3)
// 	_, err := rand.Read(b)
// 	if err != nil {
// 		return -1, err
// 	}

// 	// преобразуем байты в число
// 	userID := int(binary.BigEndian.Uint16(b))

// 	return userID, nil
// }

// // получаем ключ для генерции токена
// func getKeyForTokenGeneration() []byte {

// 	SecretKEY := "supersecretkey"

// 	return []byte(SecretKEY)
// }
