package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/cfif1982/taxi/internal/domain/routes"
)

type AddRouteDTO struct {
	Name   string `json:"name"`
	Points string `json:"points"`
}

// Обрабатываем запрос на добавление маршрута
func (h *Handler) AddRoute() http.Handler {

	fn := func(rw http.ResponseWriter, req *http.Request) {
		var addRouteDTO AddRouteDTO

		// после чтения тела запроса, закрываем
		defer req.Body.Close()

		// читаем тело запроса
		body, err := io.ReadAll(req.Body)
		if err != nil {
			h.logger.Fatal(err.Error())
		}

		if err = json.Unmarshal(body, &addRouteDTO); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		// создаем маршрут из данных запроса
		route, err := routes.CreateRoute(addRouteDTO.Name, addRouteDTO.Points)
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
