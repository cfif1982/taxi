package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/cfif1982/taxi/internal/domain/drivers"
	"github.com/cfif1982/taxi/pkg/logger"
	"github.com/google/uuid"
)

// создаю свой тип для ключей контекста. Нужно хранить id авторизованного водителя
type ctxKey string

const KeyDriverID ctxKey = "driver_id" //  ключ в контексте для поля driver_id

// Интерфейс репозитория
type DriverRepositoryInterface interface {

	// добавить водителя
	AddDriver(driver *drivers.Driver) error

	// найти водителя по телефону
	GetDriverByTelephone(telephone string) (*drivers.Driver, error)

	// найти водителя по id
	GetDriverByID(id uuid.UUID) (*drivers.Driver, error)

	// сохранить водителя
	SaveDriver(driver *drivers.Driver) error
}

// структура хэндлера
type Handler struct {
	driverRepo DriverRepositoryInterface // репозиторий для водителей
	logger     *logger.Logger            // логгер
}

// создаем новый хэндлер
func NewHandler(driverRepo DriverRepositoryInterface, logger *logger.Logger) *Handler {
	return &Handler{
		driverRepo: driverRepo,
		logger:     logger,
	}
}

// считываем данные из запроса и анмаршалим их в dto
// неможем readRequestToDTO написать как метод хэндлера, т.е. (h *Handler) readRequestToDTO
// из-за того, что в GO методы не могут иметь собственные параметры тика, как функции
// QUESTION: решил написать дженерик для считывания данных из запроса в DTO. Нормлаьный подход? этот дженерик тут нужно описать?
func readRequestToDTO[T any](h *Handler, req *http.Request, dto *T) error {

	// после чтения тела запроса, закрываем
	defer req.Body.Close()

	// читаем тело запроса
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Info(err.Error())
		return err
	}

	err = json.Unmarshal(body, dto)

	return err
}
