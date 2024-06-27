package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cfif1982/taxi/internal/domain/drivers"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/cfif1982/taxi/mocks"
	"github.com/cfif1982/taxi/pkg/logger"

	"github.com/golang/mock/gomock"
)

func TestRegistration(t *testing.T) {

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

	// Определяем тестовые случаи
	tests := []struct {
		testName           string
		requestBody        string
		mockParamDriver    *drivers.Driver
		mockParamTelephone string
		mockReturnDriver   *drivers.Driver
		mockReturnError    error
		expectedCode       int
		expectedBody       string
	}{
		{
			testName: "registration test #1",
			requestBody: `{
				"telephone": "89275656981",
				"route_id": "83d7bec6-1a15-47b8-8d58-71392f528ed7",
				"name": "cfif",
				"password": "12345"
			}`,

			mockReturnError: nil,
			expectedCode:    http.StatusOK,
		},
	}

	for _, tt := range tests {
		fmt.Println(tt.testName)

		// Настраиваем mock для текущего теста
		mockRepo.EXPECT().AddDriver(gomock.Any()).Return(tt.mockReturnError)

		// Создаем роутер и добавляем хэндлер
		routerChi := chi.NewRouter()

		handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.Registration().ServeHTTP(w, r)
		})

		routerChi.Post("/api/driver/registration", handlerFunc)

		// Создаем запрос для теста
		req, err := http.NewRequest("POST", "/api/driver/registration", bytes.NewBuffer([]byte(tt.requestBody)))
		if err != nil {
			t.Fatal(err)
		}

		// Создаем ResponseRecorder для записи ответа
		rr := httptest.NewRecorder()

		// Выполняем запрос
		routerChi.ServeHTTP(rr, req)

		// Проверяем код ответа
		assert.Equal(t, tt.expectedCode, rr.Code, "Ожидался код ответа %d", tt.expectedCode)

	}
}
