package handlers

import (
	"github.com/cfif1982/taxi/pkg/logger"
	"github.com/google/uuid"

	"github.com/cfif1982/taxi/internal/domain/routes"
)

// Интерфейс репозитория
type RouteRepositoryInterface interface {

	// Получить все маршруты
	GetAllRoutes() (*[]routes.Route, error)

	// добавить маршрут
	AddRoute(route *routes.Route) error

	// Сохранить маршрут
	SaveRoute(route *routes.Route) error

	// Удалить маршрут
	DeleteRoute(id uuid.UUID) error

	// Найти маршрут по id
	GetRouteByID(id uuid.UUID) (*routes.Route, error)
}

// структура хэндлера
type Handler struct {
	routeRepo RouteRepositoryInterface
	logger    *logger.Logger
}

// создаем новый хэндлер
func NewHandler(routeRepo RouteRepositoryInterface, logger *logger.Logger) *Handler {
	return &Handler{
		routeRepo: routeRepo,
		logger:    logger,
	}
}
