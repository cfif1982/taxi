package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/cfif1982/taxi/internal/base"
	"github.com/cfif1982/taxi/internal/domain/drivers"
	"github.com/cfif1982/taxi/pkg/logger"
)

// структура GPS
type GPS struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// переменная нужна для создания websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Обрабатываем запрос на начало работы
func (h *Handler) Start() http.HandlerFunc {

	// создаем функцию которую будем возвращать как http.HandlerFunc
	fn := func(rw http.ResponseWriter, req *http.Request) {

		// узнаем id водителя из контекста запроса
		var driverID uuid.UUID
		if req.Context().Value(KeyDriverID) != nil {
			driverID = req.Context().Value(KeyDriverID).(uuid.UUID)
		}

		// Если id водителя нет, то ошибка
		if driverID == uuid.Nil {
			http.Error(rw, drivers.ErrDriverIsNotAuthorized.Error(), http.StatusUnauthorized)
			return
		}

		// находим водителя по id
		driver, err := h.driverRepo.GetDriverByID(driverID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(rw, drivers.ErrDriverIsNotFound.Error(), http.StatusInternalServerError)
			} else {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			}

			return
		}

		// если у водителя сегодняшний день уже был оплачен, то списывания не происходит
		currentDate := time.Now().Truncate(24 * time.Hour)

		if !currentDate.Equal(driver.LastPaidDate()) {

			// узнаем стоимость работы за один день
			cost := getCost()

			// QUESTION: нужно ли здесь проверять баланс? я при списании это делаю. Ии лучше подстраховываться в таких случаях?
			if driver.Balance() < cost {
				http.Error(rw, drivers.ErrInsufficientFunds.Error(), http.StatusPaymentRequired)
				return
			}

			// списываем с баланса стоимость работы и сохраняем дату последней оплаты как сегодняшшнюю
			// QUESTION: изменение даты лучше вынести в отдельный метод?
			if err = driver.ReduceBalance(cost); err != nil {
				http.Error(rw, drivers.ErrInsufficientFunds.Error(), http.StatusPaymentRequired)
				return
			}

			// сохраняем измененные данные
			err = h.driverRepo.SaveDriver(driver)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Если всё нормально, то создаем соединение и начинаем работать
		//***************************************************************
		//***************************************************************
		// createConnection(rw, req)
		// создаем соединение websocket
		conn, err := upgrader.Upgrade(rw, req, nil)

		if err != nil {
			h.logger.Info("Ошибка подключения:", err.Error())
			return
		}
		defer conn.Close()

		h.logger.Info("Start websocket:")

		connectedDriver := base.ConnectedDriver{
			ID:                 driver.ID(),
			SendDataToDriverCH: make(chan []byte),
			DoneCH:             make(chan struct{}),
		}

		// читаем данные от водителя
		// горутина будет закрываться при закрытии websocket
		go readDataFromDriver(conn, h.logger, &connectedDriver, h.connectedDriversBase.ReceiveDataFromDriverCH)

		// здесь нужен бесконечный цикл, т.к. при завершении функции, websocket закрывается
		// но нам нужен сигнал о том, что мы действительно хотим закрыть сокет
		// поэтому ждем этот сигнал чеерез канал водителя
		// ну и до кучи чтобы еще одну горутину не делать - тут же отсылаем данные в сокет
		// данные для этого получаем через канал водителя из базы подключенных водителей
		for {
			select {
			case <-connectedDriver.DoneCH: // закрытие канала done
				h.logger.Info("Close websocket")
				return
			case data := <-connectedDriver.SendDataToDriverCH:
				// отправляем данные в сокет
				err = conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					h.logger.Info("Ошибка отправки данных:", err.Error())
					return
				}
			}

		}
	}

	return http.HandlerFunc(fn)
}

// func createConnection(rw http.ResponseWriter, req *http.Request){

// }

// читаем данные от водителя
func readDataFromDriver(
	conn *websocket.Conn,
	logger *logger.Logger,
	connectedDriver *base.ConnectedDriver,
	receiveDataFromDriverCH chan *base.ConnectedDriver) {

	// Чтение сообщений от клиента
	for {
		// Чтение сообщения
		_, message, err := conn.ReadMessage()

		if err != nil {
			// уточняем ошибку
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Info("Ошибка при чтении сообщения:", err.Error())
			}

			// Ошибка также возникает при попытке прочитать из закрытого соединения
			// Раз соединение закрыто, то завершаем горутину
			return
		}

		// структура для хранения GPS от водителя
		var gps GPS
		err = json.Unmarshal(message, &gps)

		if err != nil {
			logger.Info("Неверный формат данных от водителя", err.Error())
			break
		}

		// сохраняем данные GPS
		connectedDriver.Latitude = gps.Latitude
		connectedDriver.Longitude = gps.Longitude
		connectedDriver.ReceivedDataTime = time.Now()

		// отправляем данные в канал
		receiveDataFromDriverCH <- connectedDriver
	}
}

// Узнаем стоимость работы за один день
// Пока что это просто заглушка
func getCost() int {
	return 30
}
