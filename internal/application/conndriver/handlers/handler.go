package handlers

import (
	base "github.com/cfif1982/taxi/internal/domain/connected_drivers_base"
	"github.com/cfif1982/taxi/internal/domain/drivers"
	"github.com/cfif1982/taxi/pkg/logger"
	"github.com/google/uuid"
)

// QUESTION: у меня ctxKey и const KeyDriverID встречается в двух хэндлерах: drivers/handler.go и здесь.
// Как их лучше описать в одном месте? или оставить так?
// создаю свой тип для ключей контекста. Нужно хранить id авторизованного водителя
type ctxKey string

const KeyDriverID ctxKey = "driver_id" //  ключ в контексте для поля driver_id

// Интерфейс репозитория подключенных водителей
type ConnectedDriversBaseRepositoryInterface interface {

	// получить всех водителей из базыподключенных водителей
	GetAllDrivers() (*[]base.ConnectedDriver, error)

	// обновить данные водителя в базе подключенных водителей
	UpdateDriver(connectedDriver *base.ConnectedDriver)

	// удаляем водителя из базы подключенных водителей
	RemoveDriver(driverID uuid.UUID) error
}

// Интерфейс репозитория водителей
type DriverRepositoryInterface interface {

	// найти водителя по id
	GetDriverByID(id uuid.UUID) (*drivers.Driver, error)
}

// структура хэндлера
type Handler struct {
	driverRepo               DriverRepositoryInterface               // репозиторий водителей
	connectedDriversBaseRepo ConnectedDriversBaseRepositoryInterface // репозиторий подключенных водителей
	// connectedDriversBase *base.ConnectedDriversBase // база соединенных водителей
	logger *logger.Logger // логгер
}

// создаем новый хэндлер
func NewHandler(driverRepo DriverRepositoryInterface, connectedDriversBaseRepo ConnectedDriversBaseRepositoryInterface, logger *logger.Logger) *Handler {
	return &Handler{
		driverRepo:               driverRepo,
		connectedDriversBaseRepo: connectedDriversBaseRepo,
		logger:                   logger,
	}
}
