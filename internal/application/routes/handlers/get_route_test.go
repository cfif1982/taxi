package handlers

import (
	"bytes"
	"encoding/json"
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

func TestGetRoute(t *testing.T) {

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
		testName         string
		requestBody      string
		mockParamRouteID uuid.UUID
		mockReturnRoute  *routes.Route
		mockReturnError  error
		expectedCode     int
		expectedBody     string
	}{
		{
			testName: "get route test #1",
			requestBody: `{
  "route_id": "83d7bec6-1a15-47b8-8d58-71392f528ed7"
}`,
			mockParamRouteID: uuid.Must(uuid.Parse("83d7bec6-1a15-47b8-8d58-71392f528ed7")),
			mockReturnRoute: routes.NewRoute(
				uuid.Must(uuid.Parse("83d7bec6-1a15-47b8-8d58-71392f528ed7")),
				"31", []routes.Point{
					*routes.NewPoint(uuid.Must(uuid.Parse("83d7bec6-1a15-47b8-8d58-71392f528ed7")), "Спутник2", true, 24.34545, 10.123450),
					*routes.NewPoint(uuid.Must(uuid.Parse("83d7bec6-1a15-47b8-8d58-71392f528ed7")), "", false, 28.34545, 11.123450),
					*routes.NewPoint(uuid.Must(uuid.Parse("83d7bec6-1a15-47b8-8d58-71392f528ed7")), "Центр2", true, 26.34545, 15.123450),
				}),
			mockReturnError: nil,
			expectedCode:    http.StatusOK,
			expectedBody: `{
"id": "83d7bec6-1a15-47b8-8d58-71392f528ed7",
"name": "31",
"points": [
{
"id": "83d7bec6-1a15-47b8-8d58-71392f528ed7",
"name": "Спутник2",
"stop": true,
"latitude": 24.34545,
"longitude": 10.123450
},
{
"id": "83d7bec6-1a15-47b8-8d58-71392f528ed7",
"name": "",
"stop": false,
"latitude": 28.34545,
"longitude": 11.123450
},
{
"id": "83d7bec6-1a15-47b8-8d58-71392f528ed7",
"name": "Центр2",
"stop": true,
"latitude": 26.34545,
"longitude": 15.123450
}
]
}`,
		},
	}

	for _, tt := range tests {
		fmt.Println(tt.testName)

		// Настраиваем mock для текущего теста
		mockRepo.EXPECT().GetRouteByID(tt.mockParamRouteID).Return(tt.mockReturnRoute, tt.mockReturnError)

		// Создаем роутер и добавляем хэндлер
		routerChi := chi.NewRouter()

		handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.GetRoute().ServeHTTP(w, r)
		})

		routerChi.Get("/api/admin/route", handlerFunc)

		// Создаем запрос для теста
		req, err := http.NewRequest("GET", "/api/admin/route", bytes.NewBuffer([]byte(tt.requestBody)))
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
		// для правильного сравнения json строк парсим их в карты
		var expectedMap map[string]interface{}
		var actualMap map[string]interface{}

		// Парсинг JSON-строк в карты
		err = json.Unmarshal([]byte(tt.expectedBody), &expectedMap)
		if err != nil {
			t.Fatalf("Error unmarshaling expected JSON: %v", err)
		}

		err = json.Unmarshal([]byte(rr.Body.String()), &actualMap)
		if err != nil {
			t.Fatalf("Error unmarshaling actual JSON: %v", err)
		}

		assert.Equal(t, expectedMap, actualMap, "json тела не соответствует ожидаемому")
	}
}
