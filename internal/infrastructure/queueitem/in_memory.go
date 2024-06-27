package queueitem

import (
	"encoding/json"
	"math"
	"sync"
	"time"

	"github.com/cfif1982/taxi/internal/domain/queueitem"
	"github.com/cfif1982/taxi/pkg/logger"

	"github.com/google/uuid"
)

const SendDataPeriod = time.Second * 2 // частота отсылки данных водителю
const ConnectionPingPeriod = 5         // время задержки данных от водителя в секундах, при котором соединение считается разорванным

// DTO для отсылки данных водителю
type sendToDriverDataDTO struct {
	ID        uuid.UUID `json:"id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}

// Интерфейс обработки сообщений на сервере
type ServerMsgHandlerI interface {

	// Получение данных от водителя
	ReceiveMessageFromDriver() (*queueitem.QueueItem, error)

	// Посылаем данные на сервер
	SendMessageToServer(queueItem *queueitem.QueueItem)
}

// репозиторий очереди
type InMemoryRepository struct {
	queue            map[uuid.UUID]*queueitem.QueueItem // решил сделать map, а не слайс. В мапе быстрее будет искать по id нужного водителя
	mu               sync.Mutex                         // мьютекс для синхронизации доступа к базе
	serverMsgHandler ServerMsgHandlerI
	logger           *logger.Logger
}

func NewInMemoryRepo(serverMsgHandler ServerMsgHandlerI, logger *logger.Logger) *InMemoryRepository {
	return &InMemoryRepository{
		queue:            make(map[uuid.UUID]*queueitem.QueueItem),
		serverMsgHandler: serverMsgHandler,
		logger:           logger,
	}
}

func (r *InMemoryRepository) ServerMessageHandler() ServerMsgHandlerI {

	return r.serverMsgHandler
}

// запускаем базу в работу
func (r *InMemoryRepository) StartQueue() {

	// запускаем обработку очереди
	go r.handleQueue()

	// запускаем ожидание данных от водителя
	go r.waitDataFromDriver()
}

// запускаем обработку очереди
func (r *InMemoryRepository) handleQueue() {

	// назначаем таймер для отсылки данных водителям с нужной периодичностью
	ticker := time.NewTicker(SendDataPeriod)
	defer ticker.Stop()

	for {
		select {
		// отсылаем данные всем водителями
		case <-ticker.C:
			r.broadcastDataToAllDrivers()
		}
	}
}

// отправляем данные всем водителям
func (r *InMemoryRepository) broadcastDataToAllDrivers() {

	// QUESTION:  здесь нужно добавить defer
	r.mu.Lock()

	// формируем данные для отправки
	// слайс для отправки данных
	arrDriverDataDTO := []sendToDriverDataDTO{}

	for _, v := range r.queue {

		// сохраняем даные для формирования ответа сервера
		arrDriverDataDTO = append(
			arrDriverDataDTO,
			sendToDriverDataDTO{
				ID:        v.DriverID(),
				Latitude:  v.Latitude(),
				Longitude: v.Longitude(),
			})
	}

	// получаем строку данных для отправки
	driversString, err := json.Marshal(arrDriverDataDTO)

	if err != nil {
		r.logger.Info("Неверный формат данных в очереди:", err.Error())
		return
	}

	// пробегаемся по очереди и отправляем  данные
	// Здесь же мы определяем состояние соединения с водителем.
	// Делаем это по разнице между текущим временем и временем последнего сообщения от водителя
	for _, v := range r.queue {
		// узнаем состояние соединения с водителем
		// разница между текущим временем и временем получения последних данных от водителя
		diff := int(math.Round(time.Since(v.ReceivedDataTime()).Seconds()))

		// Если разница больше ConnectionPingPeriod, то закрываем соединение
		// Удаляем водителя из очереди. Это приведет к закрытию соединения
		if diff > ConnectionPingPeriod {
			r.RemoveDriver(v.DriverID())
			continue
		}

		// посылаю данные водителю
		v.DriverMsgHandler().SendMessageToDriver(driversString)
	}

	r.mu.Unlock()
}

// удаляем водителя из очереди
func (r *InMemoryRepository) RemoveDriver(driverID uuid.UUID) error {

	r.mu.Lock()

	// проверяем есть такой водитель в очереди
	if i, ok := r.queue[driverID]; ok {
		i.DriverMsgHandler().CloseHandler() // закрываем хэндллер
		delete(r.queue, driverID)           // удаляем водителя из очереди
	} else {
		return queueitem.ErrQueueItemNotFound
	}

	r.mu.Unlock()

	return nil
}

// получаем данные от водителей
func (r *InMemoryRepository) waitDataFromDriver() {

	for {
		// обновляем данные водителя в очереди. Если его там нет, то добавляем
		queueItem, err := r.serverMsgHandler.ReceiveMessageFromDriver()

		if err != nil {
			r.logger.Info("can't receive data from driver", err.Error())
			return
		}

		// обновляю очередь
		r.mu.Lock()

		r.queue[queueItem.DriverID()] = queueItem

		r.mu.Unlock()
	}
}
