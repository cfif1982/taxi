package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/cfif1982/taxi/internal/domain/routes"
)

// DTO для запроса и ответа
type AddRouteRequestPointDTO struct {
	Name      string  `json:"name"`
	Stop      bool    `json:"stop"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type AddRouteRequestRouteDTO struct {
	Name   string                    `json:"name"`
	Points []AddRouteRequestPointDTO `json:"points"`
}

// Обрабатываем запрос на добавление маршрута
func (h *Handler) AddRoute() http.Handler {

	fn := func(rw http.ResponseWriter, req *http.Request) {
		var routeDTO AddRouteRequestRouteDTO

		// после чтения тела запроса, закрываем
		defer req.Body.Close()

		// читаем тело запроса
		body, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		// получаем DTO из запроса
		if err = json.Unmarshal(body, &routeDTO); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		// создаем слайс точек из запроса
		points := make([]routes.Point, 0, len(routeDTO.Points))
		for _, p := range routeDTO.Points {
			points = append(points, *routes.CreatePoint(p.Name, p.Stop, p.Latitude, p.Longitude))
		}

		// создаем маршрут из данных запроса
		route := routes.CreateRoute(routeDTO.Name, points)

		// Добавляем маршрут в БД
		err = h.routeRepo.AddRoute(route)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// устанавливаем код 200
		rw.WriteHeader(http.StatusOK)

	}

	return http.HandlerFunc(fn)
}
