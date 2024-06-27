package handlers

import (
	"github.com/cfif1982/taxi/internal/domain/drivers"
	queueItemsInfra "github.com/cfif1982/taxi/internal/infrastructure/queueitem"
	"github.com/cfif1982/taxi/pkg/logger"
	"github.com/google/uuid"
)

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
