package internal

import (
	"encoding/json"
	"math"
	"time"

	"github.com/cfif1982/taxi/internal/domain/conndriver"
	"github.com/cfif1982/taxi/pkg/logger"
	"github.com/google/uuid"
)

const SendDataPeriod = time.Second * 2 // частота отсылки данных водителю
const ConnectionPingPeriod = 5         // время задержки данных от водителя при котором соединение считается разорванным

// Интерфейс репозитория подключенных водителей
type ConnectedDriversBaseRepositoryInterface interface {

	// получить всех водителей из базы подключенных водителей
	GetAllDrivers() (*[]conndriver.ConnectedDriver, error)

	// обновить данные водителя в базе подключенных водителей
	UpdateDriver(connectedDriver *conndriver.ConnectedDriver)

	// удаляем водителя из базы подключенных водителей
	RemoveDriver(driverID uuid.UUID) error
}

type sendToDriverDataDTO struct {
	ID        uuid.UUID `json:"id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}

// база подключенных к серверу водителей
type ConnectedDriversBase struct {
	connectedDriversBaseRepo ConnectedDriversBaseRepositoryInterface
	dataReceiver             conndriver.DataReceiverFromDriverInterface
	logger                   *logger.Logger
}

// Конструктор NewConnectedDriversBase
func NewConnectedDriversBase(
	connectedDriversBaseRepo ConnectedDriversBaseRepositoryInterface,
	dataReceiver conndriver.DataReceiverFromDriverInterface,
	logger *logger.Logger) *ConnectedDriversBase {
	return &ConnectedDriversBase{
		connectedDriversBaseRepo: connectedDriversBaseRepo,
		dataReceiver:             dataReceiver,
		logger:                   logger,
	}
}

// запускаем базу в работу
func (b *ConnectedDriversBase) StartBase() {

	// запускаем обработку базы подключенных водителей
	go b.handleBase()

	// запускаем получение данных от водителя
	go b.receiveDataFromDriver()
}

// обеспечиваем работу базы активных водителей
func (b *ConnectedDriversBase) handleBase() {
	// назначаем таймер для отсылки данных водителям с нужной периодичностью
	ticker := time.NewTicker(SendDataPeriod)
	defer ticker.Stop()

	for {
		select {
		// отсылаем данные всем водителями
		case <-ticker.C:
			b.broadcastDataToAllDrivers()
		}
	}
}

// получаем данные от подсоединенных водителей
func (b *ConnectedDriversBase) receiveDataFromDriver() {

	for {
		// обновляем данные водителя в базе. Если его там нет, то добавляем
		connectedDriver, err := b.dataReceiver.ReceiveDataFromDriver()

		if err != nil {
			b.logger.Info("can't receive data from driver", err.Error())
		}

		b.connectedDriversBaseRepo.UpdateDriver(connectedDriver)
	}
}

// отправляем данные всем водителям
func (b *ConnectedDriversBase) broadcastDataToAllDrivers() {

	// получаем всех водителей
	baseCopy, err := b.connectedDriversBaseRepo.GetAllDrivers()

	// слайс для отправки данных
	arrDriverDataDTO := []sendToDriverDataDTO{}

	for _, v := range *baseCopy {

		// сохраняем даные для формирования ответа сервера
		arrDriverDataDTO = append(
			arrDriverDataDTO,
			sendToDriverDataDTO{
				ID:        v.ID(),
				Latitude:  v.Latitude(),
				Longitude: v.Longitude(),
			})
	}

	// получаем строку данных для отправки
	driversString, err := json.Marshal(arrDriverDataDTO)

	if err != nil {
		b.logger.Info("Неверный формат данных в базе подключенных водителей:", err.Error())
		return
	}

	// пробегаемся по всем активным водителям и отправляем им данные
	// Здесь же мы определяем состояние соединения с водителем.
	// Делаем это по разнице между текущим временем и временем последнего сообщения от водителя
	for _, v := range *baseCopy {
		// узнаем состояние соединения с водителем
		// разница между текущим временем и временем получения последних данных от водителя
		diff := int(math.Round(time.Since(v.ReceivedDataTime()).Seconds()))

		// Если разница больше ConnectionPingPeriod, то закрываем соединение
		// Удаляем водителя из базы активных водителей. Это приведет к закрытию соединения
		if diff > ConnectionPingPeriod {
			b.removeDriverFromBase(v.ID())
			continue
		}

		v.DataSender().SendDataToDriver(driversString)
	}
}

// // обновить данные gps водителя в базе
// func (b *ConnectedDriversBase) updateDriversData(connectedDriver *ConnectedDriver) {

// 	mu.Lock()

// 	// памятка для меня: обращение к элементу map, переданной через указатель, делается через (*map_name)
// 	(*b.base)[connectedDriver.ID] = connectedDriver

// 	mu.Unlock()

// }

// удаляем водителя из базы подключенных водителей
func (b *ConnectedDriversBase) removeDriverFromBase(driverID uuid.UUID) {

	// mu.Lock()

	// // проверяем есть такой водитель в базе
	// if driver, ok := (*b.base)[driverID]; ok {
	// 	close(driver.DoneCH)             // посылаем, путем закрытия канала doneCH, сигнал о закрытии соединения
	// 	close(driver.SendDataToDriverCH) // закрываем канал передачи данных водителю
	// 	delete(*b.base, driverID)        // удаляем водителя из базы
	// } else {
	// 	b.logger.Info("Водителя нет в базе. driver_id: " + driverID.String())
	// }

	// mu.Unlock()

	close(driver.DoneCH)             // посылаем, путем закрытия канала doneCH, сигнал о закрытии соединения
	close(driver.SendDataToDriverCH) // закрываем канал передачи данных водителю
	// in memory RemoveDriver()
}
