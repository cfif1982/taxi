package routes

import (
	"github.com/google/uuid"
)

// структура для хранения точки маршрута
type Point struct {
	id        uuid.UUID
	name      string
	stop      bool // является ли эта точка остановкой транспорта
	latitude  float32
	longitude float32
}

// создаем новый объект
// нужна для использвания в других пакетах
func NewPoint(id uuid.UUID, name string, stop bool, latitude, longitude float32) *Point {
	return &Point{
		id:        id,
		name:      name,
		stop:      stop,
		latitude:  latitude,
		longitude: longitude,
	}
}

// Создаем новую точку
func CreatePoint(name string, stop bool, latitude, longitude float32) *Point {

	return NewPoint(uuid.New(), name, stop, latitude, longitude)
}

// возвращщаем поле id
func (p *Point) ID() uuid.UUID {
	return p.id
}

// возвращщаем поле name
func (p *Point) Name() string {
	return p.name
}

// возвращщаем поле stop
func (p *Point) Stop() bool {
	return p.stop
}

// возвращщаем поле latitude
func (p *Point) Latitude() float32 {
	return p.latitude
}

// возвращщаем поле longitude
func (p *Point) Longitude() float32 {
	return p.longitude
}
