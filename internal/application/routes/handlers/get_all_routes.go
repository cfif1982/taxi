package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// DTO для запроса и ответа
type GetAllRoutesResponseDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// Обрабатываем запрос на получение списка всех маршрутов. В элементах этого списка не нужен список точек маршрута,
// т.к. будут выводиться только названия маршрутов списком и всё.
func (h *Handler) GetAllRoutes() http.Handler {

	fn := func(rw http.ResponseWriter, req *http.Request) {

		arrGetAllRoutesDTO := []GetAllRoutesResponseDTO{} // слайс для хранения маршрутов

		// запрос к БД - находим все маршруты
		routes, err := h.routeRepo.GetAllRoutes()

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// сохраняем полученные данные в DTO
		for _, v := range *routes {
			arrGetAllRoutesDTO = append(
				arrGetAllRoutesDTO,
				GetAllRoutesResponseDTO{
					ID:   v.ID(),
					Name: v.Name(),
				})
		}

		// Устанавливаем в заголовке тип передаваемых данных
		rw.Header().Set("Content-Type", "application/json")

		// устанавливаем код 200
		rw.WriteHeader(http.StatusOK)

		// маршалим текст ответа
		answerText, err := json.Marshal(arrGetAllRoutesDTO)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// выводим ответ сервера
		_, err = rw.Write([]byte(answerText))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	return http.HandlerFunc(fn)
}
