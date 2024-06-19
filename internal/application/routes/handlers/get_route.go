package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
)

// DTO для запроса и ответа
type GetRouteResponsePointDTO struct {
	Name      string  `json:"name"`
	Stop      bool    `json:"stop"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type GetRouteResponseRouteDTO struct {
	ID     uuid.UUID                  `json:"route_id"`
	Name   string                     `json:"name"`
	Points []GetRouteResponsePointDTO `json:"points"`
}

type GetRouteRequestDTO struct {
	RouteID uuid.UUID `json:"route_id"`
}

// Обрабатываем запрос на получение списка всех маршрутов. В элементах этого списка не нужен список точек маршрута,
// т.к. будут выводиться только названия маршрутов списком и всё.
func (h *Handler) GetRoute() http.Handler {

	fn := func(rw http.ResponseWriter, req *http.Request) {

		var requestDTO GetRouteRequestDTO

		// после чтения тела запроса, закрываем
		defer req.Body.Close()

		// читаем тело запроса
		body, err := io.ReadAll(req.Body)
		if err != nil {
			h.logger.Info(err.Error())
		}

		if err = json.Unmarshal(body, &requestDTO); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		// запрос к БД - находим данные маршрута
		route, err := h.routeRepo.GetRouteByID(requestDTO.RouteID)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// сохраняем полученные данные в DTO
		arrPointsDTO := []GetRouteResponsePointDTO{} // слайс для хранения точек маршрута

		// QUESTION: т.к. в базе данных храним список точек в виде строки json, то получаю из БД строку точек
		// Такой же вопрос был в AddRoute. получается, что сначала я из строки запроса Unmarshal в DTO, а затем обратно часть этого DTO Marshal в строку
		// получается двойная работа. Или такой подход норм в этом случае?

		// unmarshal строку точек маршрута из БД в DTO
		if err = json.Unmarshal([]byte(route.Points()), &arrPointsDTO); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		routeDTO := GetRouteResponseRouteDTO{
			ID:     route.ID(),
			Name:   route.Name(),
			Points: arrPointsDTO,
		}

		// Устанавливаем в заголовке тип передаваемых данных
		rw.Header().Set("Content-Type", "application/json")

		// устанавливаем код 200
		rw.WriteHeader(http.StatusOK)

		// маршалим текст ответа
		answerText, err := json.Marshal(routeDTO)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// выводим ответ сервера
		_, err = rw.Write([]byte(answerText))
		if err != nil {
			h.logger.Info(err.Error())
		}
	}

	return http.HandlerFunc(fn)
}
