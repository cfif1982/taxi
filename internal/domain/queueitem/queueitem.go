package queueitem

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// список возможных шибок
var (
	ErrQueueItemNotFound = errors.New("queue item not found")
)

// Интерфейс обработки сообщений у водителя
type DriverMsgHandlerI interface {

	// Получаем данные от сервера
	ReceiveMessageFromServer() ([]byte, error)

	// Посылаем данные водителю
	SendMessageToDriver(data []byte)

	// Ждем сигнал на закрытие хэндлера
	WaitCloseSignal()

	// Закрываем хэндлер
	CloseHandler()
}

// элемент очереди, т.е. водитель с которым активно соединение
// QUESTION: а такой объект, с интерфейсом, можно будет хранить в постгрю?
type QueueItem struct {
	driverID         uuid.UUID
	latitude         float64
	longitude        float64
	receivedDataTime time.Time // время получения данных от водителя. Нужно для определения состояния соединения
	driverMsgHandler DriverMsgHandlerI
}

func NewQueueItem(
	driverID uuid.UUID,
	latitude, longitude float64,
	receivedDataTime time.Time,
	driverMsgHandler DriverMsgHandlerI) *QueueItem {
	return &QueueItem{
		driverID:         driverID,
		latitude:         latitude,
		longitude:        longitude,
		receivedDataTime: receivedDataTime,
		driverMsgHandler: driverMsgHandler,
	}
}

func (i *QueueItem) DriverID() uuid.UUID {
	return i.driverID
}

func (i *QueueItem) Latitude() float64 {
	return i.latitude
}

func (i *QueueItem) Longitude() float64 {
	return i.longitude
}

func (i *QueueItem) ReceivedDataTime() time.Time {
	return i.receivedDataTime
}

func (i *QueueItem) DriverMsgHandler() DriverMsgHandlerI {
	return i.driverMsgHandler
}
