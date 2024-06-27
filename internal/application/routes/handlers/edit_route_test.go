package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cfif1982/taxi/internal/domain/routes"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/cfif1982/taxi/mocks"
	"github.com/cfif1982/taxi/pkg/logger"

	"github.com/golang/mock/gomock"
)

func TestEditRoute(t *testing.T) {

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
		mockParamRoute  *routes.Route
		mockReturnError error
		expectedCode    int
	}{
		{
			testName: "edit test #1",
			requestBody: `{
"id": "83d7bec6-1a15-47b8-8d58-71392f528ed7",
"name": "31",
"points": [
{
"id": "83d7bec6-1a15-47b8-8d58-71392f528ed7",
"name": "Спутник2",
"stop": true,
"latitude": 24.3454523,
"longitude": 10.123450
},
{
"id": "83d7bec6-1a15-47b8-8d58-71392f528ed7",
"name": "",
"stop": false,
"latitude": 28.3454523,
"longitude": 11.123450
},
{
"id": "83d7bec6-1a15-47b8-8d58-71392f528ed7",
"name": "Центр2",
"stop": true,
"latitude": 26.3454523,
"longitude": 15.123450
}
]
}`,
			mockParamRoute: routes.NewRoute(
				uuid.Must(uuid.Parse("83d7bec6-1a15-47b8-8d58-71392f528ed7")),
				"31", []routes.Point{
					*routes.NewPoint(uuid.Must(uuid.Parse("83d7bec6-1a15-47b8-8d58-71392f528ed7")), "Спутник2", true, 24.3454523, 10.123450),
					*routes.NewPoint(uuid.Must(uuid.Parse("83d7bec6-1a15-47b8-8d58-71392f528ed7")), "", false, 28.3454523, 11.123450),
					*routes.NewPoint(uuid.Must(uuid.Parse("83d7bec6-1a15-47b8-8d58-71392f528ed7")), "Центр2", true, 26.3454523, 15.123450),
				}),
			mockReturnError: nil,
			expectedCode:    http.StatusOK,
		},
	}

	for _, tt := range tests {
		fmt.Println(tt.testName)

		// Настраиваем mock для текущего теста
		mockRepo.EXPECT().SaveRoute(tt.mockParamRoute).Return(tt.mockReturnError)

		// Создаем роутер и добавляем хэндлер
		routerChi := chi.NewRouter()

		handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.EditRoute().ServeHTTP(w, r)
		})

		routerChi.Put("/api/admin/route", handlerFunc)

		// Создаем запрос для теста
		req, err := http.NewRequest("PUT", "/api/admin/route", bytes.NewBuffer([]byte(tt.requestBody)))
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
