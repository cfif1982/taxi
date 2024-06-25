package handlers

import (
	"github.com/cfif1982/taxi/internal/domain/drivers"
	queueItemsInfra "github.com/cfif1982/taxi/internal/infrastructure/queueitem"
	"github.com/cfif1982/taxi/pkg/logger"
	"github.com/google/uuid"
)

// // QUESTION: у меня ctxKey и const KeyDriverID встречается в двух хэндлерах: drivers/handler.go и здесь.
// // Как их лучше описать в одном месте? или оставить так?
// // создаю свой тип для ключей контекста. Нужно хранить id авторизованного водителя
// type ctxKey string

// const KeyDriverID ctxKey = "driver_id" //  ключ в контексте для поля driver_id

// Интерфейс репозитория водителей
type DriverRepositoryInterface interface {

	// найти водителя по id
	GetDriverByID(id uuid.UUID) (*drivers.Driver, error)

	// сохранить водителя
	SaveDriver(driver *drivers.Driver) error
}

// структура хэндлера
type Handler struct {
	driverRepo           DriverRepositoryInterface // репозиторий водителей
	serverMessageHandler queueItemsInfra.ServerMsgHandlerI
	logger               *logger.Logger // логгер
}

// создаем новый хэндлер
func NewHandler(driverRepo DriverRepositoryInterface, serverMessageHandler queueItemsInfra.ServerMsgHandlerI, logger *logger.Logger) *Handler {
	return &Handler{
		driverRepo:           driverRepo,
		serverMessageHandler: serverMessageHandler,
		logger:               logger,
	}
}
