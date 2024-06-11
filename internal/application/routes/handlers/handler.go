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

	// Редактировать маршрут
	EditRoute(route *routes.Route) error

	// Удалить маршрут
	DeleteRoute(route *routes.Route) error

	// Найти маршрут по id
	GetRouteById(id uuid.UUID) (*routes.Route, error)
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