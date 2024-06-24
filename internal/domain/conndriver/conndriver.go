package conndriver

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// список возможных шибок
var (
	ErrDriverNotFound = errors.New("driver not found")
)

// подключенный к серверу водитель, т.е. водитель с которым активно соединение
// QUESTION: а такой объект, с интерфейсом, можно будет хранить в постгрю?
type ConnectedDriver struct {
	id               uuid.UUID
	latitude         float64
	longitude        float64
	receivedDataTime time.Time // время получения данных от водителя. Нужно для определения состояния соединения
	dataSender       DataSenderToDriverInterface
}

func NewConnectedDriver(id uuid.UUID, latitude, longitude float64, receivedDataTime time.Time) *ConnectedDriver {
	return &ConnectedDriver{
		id:               id,
		latitude:         latitude,
		longitude:        longitude,
		receivedDataTime: receivedDataTime,
	}
}

func (d *ConnectedDriver) ID() uuid.UUID {
	return d.id
}

func (d *ConnectedDriver) Latitude() float64 {
	return d.latitude
}

func (d *ConnectedDriver) Longitude() float64 {
	return d.longitude
}

func (d *ConnectedDriver) ReceivedDataTime() time.Time {
	return d.receivedDataTime
}

func (d *ConnectedDriver) DataSender() DataSenderToDriverInterface {
	return d.dataSender
}
