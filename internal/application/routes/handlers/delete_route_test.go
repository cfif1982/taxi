package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/cfif1982/taxi/mocks"
	"github.com/cfif1982/taxi/pkg/logger"

	"github.com/golang/mock/gomock"
)

func TestDeleteRoute(t *testing.T) {

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
		mockReturnError  error
		expectedCode     int
	}{
		{
			testName: "delete route test #1",
			requestBody: `{
  "route_id": "83d7bec6-1a15-47b8-8d58-71392f528ed7"
}`,
			mockParamRouteID: uuid.Must(uuid.Parse("83d7bec6-1a15-47b8-8d58-71392f528ed7")),
			mockReturnError:  nil,
			expectedCode:     http.StatusOK,
		},
	}

	for _, tt := range tests {
		fmt.Println(tt.testName)

		// Настраиваем mock для текущего теста
		mockRepo.EXPECT().DeleteRoute(tt.mockParamRouteID).Return(tt.mockReturnError)

		// Создаем роутер и добавляем хэндлер
		routerChi := chi.NewRouter()

		handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.DeleteRoute().ServeHTTP(w, r)
		})

		routerChi.Delete("/api/admin/route", handlerFunc)

		// Создаем запрос для теста
		req, err := http.NewRequest("DELETE", "/api/admin/route", bytes.NewBuffer([]byte(tt.requestBody)))
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
