package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/cfif1982/taxi/internal/application"
	"github.com/cfif1982/taxi/internal/domain/drivers"
	"github.com/google/uuid"
)

// DTO для ответа
type GetBalanceResponseDTO struct {
	Balance int `json:"balance"`
}

// Обрабатываем запрос на авторизацию водителя
func (h *Handler) GetBalance() http.HandlerFunc {

	// создаем функцию которую будем возвращать как http.HandlerFunc
	fn := func(rw http.ResponseWriter, req *http.Request) {

		// узнаем id водителя из контекста запроса
		var driverID uuid.UUID
		if req.Context().Value(application.KeyDriverID) != nil {
			driverID = req.Context().Value(application.KeyDriverID).(uuid.UUID)
		}

		// Если id водителя нет, то ошибка
		if driverID == uuid.Nil {
			http.Error(rw, drivers.ErrDriverIsNotAuthorized.Error(), http.StatusUnauthorized)
			return
		}

		// находим водителя по id
		driver, err := h.driverRepo.GetDriverByID(driverID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(rw, drivers.ErrDriverIsNotFound.Error(), http.StatusInternalServerError)
			} else {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			}

			return
		}

		balanceDTO := GetBalanceResponseDTO{
			Balance: driver.Balance(),
		}

		// Устанавливаем в заголовке тип передаваемых данных
		rw.Header().Set("Content-Type", "application/json")

		// устанавливаем код 200
		rw.WriteHeader(http.StatusOK)

		// маршалим текст ответа
		answerText, err := json.Marshal(balanceDTO)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// выводим ответ сервера
		_, err = rw.Write([]byte(answerText))
		if err != nil {
			h.logger.Info(err.Error())
		}
	}

	return http.HandlerFunc(fn)
}
