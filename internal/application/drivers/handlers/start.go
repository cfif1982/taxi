package handlers

import (
	"database/sql"
	"encoding/json"
	"math"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/cfif1982/taxi/internal/base"
	"github.com/cfif1982/taxi/internal/domain/drivers"
	"github.com/cfif1982/taxi/pkg/logger"
)

const SendDataPeriod = time.Second * 2 // частота отсылки данных водителю
const ConnectionPingPeriod = 5         // время задержки данных от водителя при котором соединение считается разорванным

// // DTO для получения координат от водителя
// type gpsDTO struct {
// 	Latitude  float64 `json:"latitude"`
// 	Longitude float64 `json:"longitude"`
// }

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

		receiveFromDriverCH := make(chan []byte) // Создаем канал для получения данные от водителя
		done := make(chan bool)                  // Создаем сигнальный канал для закрытия websocket
		var lastMessageReceivedTime time.Time    // время получения последнего сообщения от водителя. Нужно для проверки состояния соединения

		// var driverGPS gpsDTO // храним полученные от водителя координаты

		// создаем соединение websocket
		conn, err := upgrader.Upgrade(rw, req, nil)

		if err != nil {
			h.logger.Info("Ошибка подключения:", err.Error())
			return
		}
		defer conn.Close()

		h.logger.Info("Start websocket:")

		lastMessageReceivedTime = time.Now() // перед началом работы инициализируем время получения сообщения от водителя

		// читаем данные от водителя
		// горутина будет закрываться при закрытии websocket
		go readDataFromDriver(conn, h.logger, receiveFromDriverCH)

		// Посылаем данные водителю с заданной периодичностью
		// Канал done закрывается по истечении ConnectionPingPeriod секунд.
		// Этот канал нужен для завершения цикла ожидания данных из горутины readDataFromDriver
		// В этом случае и сама горутина закрывается
		go sendDataToDriver(conn, h.activeDriversBase, h.logger, done, &lastMessageReceivedTime)

		// горутина проверки состояния соединения
		// go checkConnectionState(&lastMessageReceivedTime, h.logger, done)

		// Ожидаем результат из канала от водителя
		for {
			select {
			case <-done: // закрытие канала done
				h.logger.Info("Close websocket")
				return
			case receivedDataFromDriver, ok := <-receiveFromDriverCH: // Ожидаем результат из канала от водителя
				if !ok {
					h.logger.Info("Channel receiveFromDriverCH closed")
					return
				}

				lastMessageReceivedTime = time.Now()

				h.activeDriversBase.UpdateGPSData(receivedDataFromDriver)

				// err = json.Unmarshal(receivedDataFromDriver, &driverGPS)

				// // если получили правильные данные, то изменяем их в массиве активных водителей
				// if err != nil {
				// 	h.logger.Info("Неверный формат данных GPS от водителя:", err.Error())
				// } else {
				// 	// меняем GPS данные у водителя
				// 	err = driver.ChangeGPS(driverGPS.Latitude, driverGPS.Longitude)

				// 	if err == nil {
				// 		// обновляем водителя в мапе
				// 		h.activeDriversBase.UpdateGPSData(receivedDataFromDriver)
				// 	}
				// }
			}
		}
	}

	return http.HandlerFunc(fn)
}

// Посылаем данные водителю с заданной периодичностью
func sendDataToDriver(
	conn *websocket.Conn,
	arrActiveDrivers *base.ActiveDriversBase,
	logger *logger.Logger,
	done chan bool,
	lastMessageTime *time.Time) {

	ticker := time.NewTicker(SendDataPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			// получаем разницу между временем последнего сообщения от водителя
			// эту разницу в секундах округляем до целого значения
			diff := int(math.Round(time.Now().Sub(*lastMessageTime).Seconds()))

			// Если разница больше ConnectionPingPeriod, то закрываем соединение
			// Закрываем канал done. Это приведет к завершению цикла ожидания сообщений от водителя
			// и там же закроется соединение
			if diff > ConnectionPingPeriod {
				close(done)
				return
			}

			logger.Info("тик таймера")
			sendData := []sendDataToDriverDTO{} // слайс для хранения отправляемых данных

			// сохраняем отправляемые данные в DTO
			for _, v := range *arrActiveDrivers {
				sendData = append(
					sendData,
					sendDataToDriverDTO{
						DriverID:  v.ID(),
						Latitude:  v.Latitude(),
						Longitude: v.Longitude(),
					})
			}

			// маршалим отправляемый текст
			sendText, err := json.Marshal(sendData)

			// отправляем текст
			err = conn.WriteMessage(websocket.TextMessage, []byte(sendText))
			if err != nil {
				logger.Info("Ошибка отправки данных:", err.Error())
				return
			}
		}
	}
}

// читаем данные от водителя
func readDataFromDriver(
	conn *websocket.Conn,
	logger *logger.Logger,
	ch chan []byte) {

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

		// отправляем данные в канал
		ch <- message
	}
}

// Узнаем стоимость работы за один день
// Пока что это просто заглушка
func getCost() int {
	return 30
}
