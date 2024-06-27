package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/cfif1982/taxi/mocks"
	"github.com/cfif1982/taxi/pkg/logger"

	"github.com/golang/mock/gomock"
)

func TestAddRoute(t *testing.T) {

	// создаём контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	mockRepo := mocks.NewMockRouteRepositoryInterface(ctrl)

	// инициализируем логгер
	logger, err := logger.GetLogger()

	// Если логгер не инициализировался, то выводим сообщение с помощью обычного log
	if err != nil {
		log.Fatal("logger zap initialization: FAILURE")
	}

	handler := NewHandler(mockRepo, logger)

	// Определяем тестовые случаи
	tests := []struct {
		testName        string
		requestBody     string
		mockReturnError error
		expectedCode    int
	}{
		{
			testName: "add route test #1",
			requestBody: `{
"name": "29",
"points": [
{
"name": "Спутник",
"stop": true,
"latitude": 24.3454523,
"longitude": 10.123450
},
{
"name": "",
"stop": false,
"latitude": 28.3454523,
"longitude": 11.123450
},
{
"name": "Центр",
"stop": true,
"latitude": 26.3454523,
"longitude": 15.123450
}
]
}`,
			mockReturnError: nil,
			expectedCode:    http.StatusOK,
		},
	}

	for _, tt := range tests {
		fmt.Println(tt.testName)

		// Настраиваем mock для текущего теста
		mockRepo.EXPECT().AddRoute(gomock.Any()).Return(tt.mockReturnError)

		// Создаем роутер и добавляем хэндлер
		routerChi := chi.NewRouter()

		handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.AddRoute().ServeHTTP(w, r)
		})

		routerChi.Post("/api/admin/route", handlerFunc)

		// Создаем запрос для теста
		req, err := http.NewRequest("POST", "/api/admin/route", bytes.NewBuffer([]byte(tt.requestBody)))
		if err != nil {
			t.Fatal(err)
		}

		// Создаем ResponseRecorder для записи ответа
		rr := httptest.NewRecorder()

		// Выполняем запрос
		routerChi.ServeHTTP(rr, req)

		// Проверяем код ответа
		assert.Equal(t, tt.expectedCode, rr.Code, "Ожидался код ответа %d", tt.expectedCode)

		// Проверяем тело ответа
		// assert.Equal(t, tt.expectedBody, rr.Body.String(), "Ответ тела не соответствует ожидаемому для userID %s", tt.userID)
	}
}
