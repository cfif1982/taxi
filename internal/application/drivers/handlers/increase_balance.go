package handlers

import (
	"database/sql"
	"net/http"

	"github.com/cfif1982/taxi/internal/domain/drivers"
)

type DriverIncreaseBalanceRequest struct {
	Telephone string `json:"telephone"`
	Summa     int    `json:"summa"`
}

// Обрабатываем запрос на авторизацию водителя
func (h *Handler) IncreaseBalance() http.HandlerFunc {

	// создаем функцию которую будем возвращать как http.HandlerFunc
	fn := func(rw http.ResponseWriter, req *http.Request) {

		var driverIncBalance DriverIncreaseBalanceRequest

		// считываем запрос в dto
		if err := readRequestToDTO(h, req, &driverIncBalance); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		// нахоим водителя по телефону
		driver, err := h.driverRepo.GetDriverByTelephone(driverIncBalance.Telephone)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(rw, drivers.ErrWrongPassword.Error(), http.StatusUnauthorized)
			} else {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			}

			return
		}

		// увеличиваем баланс
		err = driver.IncreaseBalance(driverIncBalance.Summa)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// сохраняем измененные данные
		err = h.driverRepo.SaveDriver(driver)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// устанавливаем код 200
		rw.WriteHeader(http.StatusOK)
	}

	return http.HandlerFunc(fn)
}
