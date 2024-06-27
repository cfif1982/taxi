package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/cfif1982/taxi/internal/application"

	"github.com/cfif1982/taxi/internal/domain/drivers"
	"github.com/cfif1982/taxi/internal/domain/queueitem"
	queueItemsInfra "github.com/cfif1982/taxi/internal/infrastructure/queueitem"
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
		if req.Context().Value(application.KeyDriverID) != nil {
			driverID = req.Context().Value(application.KeyDriverID).(uuid.UUID)
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

			// списываем с баланса стоимость работы и сохраняем дату последней оплаты как сегодняшшнюю
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
		// создаем соединение websocket
		conn, err := upgrader.Upgrade(rw, req, nil)

		if err != nil {
			h.logger.Info("Ошибка подключения:", err.Error())
			return
		}
		defer conn.Close()

		h.logger.Info("Start websocket:")

		// Создаем обработчик сообщений водителя - DriverMsgHandler
		// DriverMsgHandler отвечает за прием сообщений от сервера
		// сейчас он реализован через каналы. Но можно реализовать и через RabbitMQ
		driverMsgHandler := queueItemsInfra.NewChannelDriverMsgHandler()

		// ожидаем данные из сокета
		// горутина будет закрываться при закрытии websocket
		go waitDataFromSocket(conn, driverID, driverMsgHandler, h.serverMessageHandler, h.logger)

		// ожидаем данные от сервера
		// горутина будет закрываться при получении сигнала на закрытие
		go waitDataFromServer(conn, driverMsgHandler, h.logger)

		// здесь нужен бесконечный цикл, т.к. при завершении функции, websocket закрывается
		// но нам нужен сигнал о том, что мы действительно хотим закрыть сокет
		// поэтому ждем этот сигнал чеерез хэндлер  водителя
		driverMsgHandler.WaitCloseSignal()
	}

	return http.HandlerFunc(fn)
}

// ожидаем данные от сервера
func waitDataFromServer(
	conn *websocket.Conn,
	driverMsgHandler queueitem.DriverMsgHandlerI,
	logger *logger.Logger) {

	for {
		// ждем данные от сервера
		driversString, err := driverMsgHandler.ReceiveMessageFromServer()

		// если канал закрывается, то произойдет ошибка
		if err != nil {
			if err == queueItemsInfra.ErrDriverChannelClosed {
				logger.Info("channel closed")
				return
			}

			logger.Info("error:", err.Error())
			return
		}

		// отправляем данные в сокет
		err = conn.WriteMessage(websocket.TextMessage, driversString)
		if err != nil {
			logger.Info("Ошибка отправки данных:", err.Error())
			return
		}
	}
}

// ожидаем данные из сокета
func waitDataFromSocket(
	conn *websocket.Conn,
	driverID uuid.UUID,
	driverMsgHandler queueitem.DriverMsgHandlerI,
	serverMsgHandler queueItemsInfra.ServerMsgHandlerI,
	logger *logger.Logger) {

	// Чтение сообщений от водителя
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

		// создаем элемент очереди
		queueItem := queueitem.NewQueueItem(
			driverID,
			gps.Latitude,
			gps.Longitude,
			time.Now(),
			driverMsgHandler)

		// отправляем данные в канал
		serverMsgHandler.SendMessageToServer(queueItem)
	}
}

// Узнаем стоимость работы за один день
// Пока что это просто заглушка
func getCost() int {
	return 30
}
