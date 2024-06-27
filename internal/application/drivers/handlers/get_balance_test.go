package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cfif1982/taxi/internal/application"
	"github.com/cfif1982/taxi/internal/domain/drivers"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/cfif1982/taxi/mocks"
	"github.com/cfif1982/taxi/pkg/logger"

	"github.com/golang/mock/gomock"
)

func TestGetBalance(t *testing.T) {

	// создаём контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	mockRepo := mocks.NewMockDriverRepositoryInterface(ctrl)

	// инициализируем логгер
	logger, err := logger.GetLogger()

	// Если логгер не инициализировался, то выводим сообщение с помощью обычного log
	if err != nil {
		log.Fatal("logger zap initialization: FAILURE")
	}

	handler := NewHandler(mockRepo, logger)

	var zeroTime time.Time

	// Определяем тестовые случаи
	tests := []struct {
		testName           string
		requestBody        string
		mockParamTelephone string
		mockReturnDriver   *drivers.Driver
		mockReturnError    error
		expectedCode       int
		expectedBody       string
		ctxParamDriver     uuid.UUID
	}{
		{
			testName: "get balance test #1",
			mockReturnDriver: drivers.NewDriver(
				uuid.Must(uuid.Parse("83d7bec6-1a15-47b8-8d58-71392f528ed7")),
				uuid.Must(uuid.Parse("83d7bec6-1a15-47b8-8d58-71392f528ed7")),
				"89275656981", "cfif", "12345", 300, zeroTime),
			mockReturnError: nil,
			expectedCode:    http.StatusOK,
			expectedBody: `{
				"balance": 300}`,
			ctxParamDriver: uuid.Must(uuid.Parse("83d7bec6-1a15-47b8-8d58-71392f528ed7")),
		},
	}

	for _, tt := range tests {
		fmt.Println(tt.testName)

		// Настраиваем mock для текущего теста
		mockRepo.EXPECT().GetDriverByID(gomock.Any()).Return(tt.mockReturnDriver, tt.mockReturnError)

		// Создаем роутер и добавляем хэндлер
		routerChi := chi.NewRouter()

		handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.GetBalance().ServeHTTP(w, r)
		})

		routerChi.Get("/api/driver/balance", handlerFunc)

		// Создаем запрос для теста
		req, err := http.NewRequest("GET", "/api/driver/balance", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Создание базового контекста
		ctx := context.Background()

		// создаю контекст для сохранения userID
		ctx = context.WithValue(ctx, application.KeyDriverID, tt.ctxParamDriver)

		// Добавление контекста в запрос
		req = req.WithContext(ctx)

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
