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

		if err = json.Unmarshal(body, &routeDTO); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		// QUESTION: т.к. в базе данных храним список точек в виде строки json, то получаю обратно строку точек
		// правильно делаю? получается, что сначала я из строки запроса Unmarshal в DTO, а затем обратно часть этого DTO Marshal в строку
		// получается двойная работа. Или такой подход норм в этом случае?
		routePointsString, _ := json.Marshal(routeDTO.Points)

		// создаем маршрут из данных запроса
		route, err := routes.CreateRoute(routeDTO.Name, string(routePointsString))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

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
