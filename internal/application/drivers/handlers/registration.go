package handlers

import (
	"net/http"

	"github.com/cfif1982/taxi/internal/domain/drivers"
	"github.com/google/uuid"
)

// DTO для запроса и ответа
type RegistrationRequestDTO struct {
	RouteID   uuid.UUID `json:"route_id"`
	Telephone string    `json:"telephone"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
}

// Обрабатываем запрос на регистрацию водителя
func (h *Handler) Registration() http.HandlerFunc {

	fn := func(rw http.ResponseWriter, req *http.Request) {

		var regDTO RegistrationRequestDTO

		// считываем запрос в dto
		if err := readRequestToDTO(h, req, &regDTO); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		// создаем водителя из данных запроса
		driver, err := drivers.CreateDriver(regDTO.RouteID, regDTO.Telephone, regDTO.Name, regDTO.Password)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// Добавляем водителя в БД
		err = h.driverRepo.AddDriver(driver)
		if err != nil {
			// Если телефон уже занят
			if err == drivers.ErrTelephoneAlreadyExist {
				http.Error(rw, err.Error(), http.StatusConflict)
			} else {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// устанавливаем код 200
		rw.WriteHeader(http.StatusOK)

	}

	return http.HandlerFunc(fn)
}
