package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type DeleteRouteRequestDTO struct {
	RouteID uuid.UUID `json:"route_id"`
}

// Обрабатываем запрос на получение списка всех маршрутов. В элементах этого списка не нужен список точек маршрута,
// т.к. будут выводиться только названия маршрутов списком и всё.
func (h *Handler) DeleteRoute() http.Handler {

	fn := func(rw http.ResponseWriter, req *http.Request) {

		var requestDTO DeleteRouteRequestDTO

		// после чтения тела запроса, закрываем
		defer req.Body.Close()

		// читаем тело запроса
		body, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = json.Unmarshal(body, &requestDTO); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		// запрос к БД - удаляем маршрут
		err = h.routeRepo.DeleteRoute(requestDTO.RouteID)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// устанавливаем код 200
		rw.WriteHeader(http.StatusOK)

	}

	return http.HandlerFunc(fn)
}
