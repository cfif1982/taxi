package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/cfif1982/taxi/internal/domain/routes"
	"github.com/google/uuid"
)

// DTO для запроса и ответа
type EditRouteRequestPointDTO struct {
	Name      string  `json:"name"`
	Stop      bool    `json:"stop"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type EditRouteRequestRouteDTO struct {
	ID     uuid.UUID                  `json:"route_id"`
	Name   string                     `json:"name"`
	Points []EditRouteRequestPointDTO `json:"points"`
}

// Обрабатываем запрос на обновление данных маршрута
func (h *Handler) EditRoute() http.Handler {

	fn := func(rw http.ResponseWriter, req *http.Request) {

		var routeDTO EditRouteRequestRouteDTO

		// после чтения тела запроса, закрываем
		defer req.Body.Close()

		// читаем тело запроса
		body, err := io.ReadAll(req.Body)
		if err != nil {
			h.logger.Fatal(err.Error())
		}

		// получаем DTO из запроса
		if err = json.Unmarshal(body, &routeDTO); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		// Получаем строку точек маршрута из DTO запроса.
		// QUESTION: здесь всё тот же вопрос про маршал анмаршал - двойная работа.
		routePointsString, _ := json.Marshal(routeDTO.Points)

		// создаем объект маршрута по данным из запроса
		route, err := routes.NewRoute(routeDTO.ID, routeDTO.Name, string(routePointsString))

		// запрос к БД - сохраняем измененные данные маршрута
		err = h.routeRepo.EditRoute(route)

		if err != nil {
			h.logger.Info(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// устанавливаем код 200
		rw.WriteHeader(http.StatusOK)

	}

	return http.HandlerFunc(fn)
}
