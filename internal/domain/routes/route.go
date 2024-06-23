package routes

import (
	"errors"

	"github.com/google/uuid"
)

// список возможных шибок
var (
	ErrRouteNotFound    = errors.New("route not found")
	ErrNameAlreadyExist = errors.New("route with this name already exist")
)

// структура для хранения маршрута
type Route struct {
	id     uuid.UUID
	name   string
	points string
}

// создаем новый объект
// нужна для использвания в других пакетах
func NewRoute(id uuid.UUID, name, points string) (*Route, error) {
	return &Route{
		id:     id,
		name:   name,
		points: points,
	}, nil
}

// Создаем новый маршрут
func CreateRoute(name, points string) (*Route, error) {

	return NewRoute(uuid.New(), name, points)

}

// возвращщаем поле id
func (r *Route) ID() uuid.UUID {
	return r.id
}

// возвращщаем поле name
func (r *Route) Name() string {
	return r.name
}

// возвращщаем поле points
func (r *Route) Points() string {
	return r.points
}
