package conndriver

import (
	"encoding/json"
	"math"
	"time"

	"github.com/google/uuid"
)

const ConnectionPingPeriod = 5 // время задержки данных от водителя при котором соединение считается разорванным

// Создаем базу
func CreateConnectedDriversBase(dataReceiver *ReceiverDataFromDriver) *ConnectedDriversBase {

	return &ConnectedDriversBase{
		dataReceiver: dataReceiver,
	}
}

type sendToDriverDataDTO struct {
	ID        uuid.UUID `json:"id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}

// отправляем данные всем водителям
func (b *ConnectedDriversBase) broadcastDataToAllDrivers() {

	// mu.Lock()

	// // делаем копию базы, чтобы можно было по ней пробежаться не мешая добавлению данных
	// baseCopy := *b.base

	// mu.Unlock()

	// вместо этого используем
	// func (r *InMemoryRepository) GetAllDrivers() (*[]base.ConnectedDriver, error)
	// baseCopy = in_memory.GetAllDrivers()

	// слайс для отправки данных
	arrDriverDataDTO := []sendToDriverDataDTO{}

	for _, v := range baseCopy {

		// сохраняем даные для формирования ответа сервера
		arrDriverDataDTO = append(
			arrDriverDataDTO,
			sendToDriverDataDTO{
				ID:        v.ID,
				Latitude:  v.Latitude,
				Longitude: v.Longitude,
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
	for _, v := range baseCopy {
		// узнаем состояние соединения с водителем
		// разница между текущим временем и временем получения последних данных от водителя
		diff := int(math.Round(time.Since(v.ReceivedDataTime).Seconds()))

		// Если разница больше ConnectionPingPeriod, то закрываем соединение
		// Удаляем водителя из базы активных водителей. Это приведет к закрытию соединения
		if diff > ConnectionPingPeriod {
			b.removeDriverFromBase(v.ID)
			continue
		}

		b.logger.Info("send string: " + string(driversString))
		v.SendDataToDriverCH <- driversString
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
