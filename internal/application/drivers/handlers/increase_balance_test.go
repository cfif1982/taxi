package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cfif1982/taxi/internal/domain/drivers"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/cfif1982/taxi/mocks"
	"github.com/cfif1982/taxi/pkg/logger"

	"github.com/golang/mock/gomock"
)

func TestIncreaseBalance(t *testing.T) {

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
	}{
		{
			testName: "increase balance test #1",
			requestBody: `{
				"telephone": "89275656981",
				"summa": 200
			}`,
			mockParamTelephone: "89275656981",
			mockReturnDriver: drivers.NewDriver(
				uuid.Must(uuid.Parse("83d7bec6-1a15-47b8-8d58-71392f528ed7")),
				uuid.Must(uuid.Parse("83d7bec6-1a15-47b8-8d58-71392f528ed7")),
				"89275656981", "cfif", "12345", 300, zeroTime),
			mockReturnError: nil,
			expectedCode:    http.StatusOK,
		},
	}

	for _, tt := range tests {
		fmt.Println(tt.testName)

		// Настраиваем mock для текущего теста
		mockRepo.EXPECT().GetDriverByTelephone(tt.mockParamTelephone).Return(tt.mockReturnDriver, tt.mockReturnError)
		mockRepo.EXPECT().SaveDriver(gomock.Any()).Return(tt.mockReturnError)

		// Создаем роутер и добавляем хэндлер
		routerChi := chi.NewRouter()

		handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.IncreaseBalance().ServeHTTP(w, r)
		})

		routerChi.Put("/api/driver/balance", handlerFunc)

		// Создаем запрос для теста
		req, err := http.NewRequest("PUT", "/api/driver/balance", bytes.NewBuffer([]byte(tt.requestBody)))
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
