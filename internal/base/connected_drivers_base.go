package base

import (
	"encoding/json"
	"math"
	"sync"
	"time"

	"github.com/cfif1982/taxi/pkg/logger"
	"github.com/google/uuid"
)

const ConnectionPingPeriod = 5         // время задержки данных от водителя при котором соединение считается разорванным
const SendDataPeriod = time.Second * 2 // частота отсылки данных водителю

var mu sync.Mutex // мьютекс для синхронизации доступа к базе

type sendToDriverDataDTO struct {
	ID        uuid.UUID `json:"id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}

// подключенный к серверу водитель, т.е. водитель с которым активно соединение websocket
type ConnectedDriver struct {
	ID                 uuid.UUID     `json:"id"`
	Latitude           float64       `json:"latitude"`
	Longitude          float64       `json:"longitude"`
	ReceivedDataTime   time.Time     // время получения данных от водителя. Нужно для определения состояния соединения
	SendDataToDriverCH chan []byte   // канал, через который буду отсылаться данные водителю
	DoneCH             chan struct{} // канал, по которому будет передаваться сигнал о закрытии горутин приема и отправки данных
}

// структура базы подключенных водителей
type ConnectedDriversBase struct {
	base                    *map[uuid.UUID]*ConnectedDriver // решил сделать map, а не слайс. В мапе быстрее будет искать по id нужного водителя
	logger                  *logger.Logger                  // логгер
	receiveDataFromDriverCH chan *ConnectedDriver           // канал, по которому будут передаваться данные от водителей
}

// Создаем базу
func CreateConnectedDriversBase(logger *logger.Logger) (*ConnectedDriversBase, error) {

	var base = make(map[uuid.UUID]*ConnectedDriver)

	return &ConnectedDriversBase{
		base:                    &base,
		logger:                  logger,
		receiveDataFromDriverCH: make(chan *ConnectedDriver),
		// RemoveDriverCH:          make(chan uuid.UUID),
	}, nil
}

// обеспечиваем работу базы активных водителей
func (b *ConnectedDriversBase) HandleBase() {
	// назначаем таймер для отсылки данных водителям с нужной периодичностью
	ticker := time.NewTicker(SendDataPeriod)
	defer ticker.Stop()

	for {
		select {
		// отсылаем данные всем водителями
		case <-ticker.C:
			b.broadcastDataToAllDrivers()

		// обновляем данные водителя в базе. Если его там нет, то добавляем
		case connectedDriver := <-b.receiveDataFromDriverCH:
			b.updateDriversData(connectedDriver)

			// удаляем неактивного водителя из базы. Пока не знаю - нужно это или нет
			// водитель хочет завершить работу. Жмет у себя кнопку закрыть приложение (ну или отключиться) - вот в этом слугае от него и идет такой сигнал.
			// case driverID := <-b.RemoveDriverCH:
			// 	b.removeDriverFromBase(driverID)
		}
	}
}

// отправляем данные всем водителям
func (b *ConnectedDriversBase) broadcastDataToAllDrivers() {

	mu.Lock()

	// делаем копию базы, чтобы можно было по ней пробежаться не мешая добавлению данных
	baseCopy := *b.base

	mu.Unlock()

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

// обновить данные gps водителя в базе
func (b *ConnectedDriversBase) updateDriversData(connectedDriver *ConnectedDriver) {

	mu.Lock()

	// памятка для меня: обращение к элементу map, переданной через указатель, делается через (*map_name)
	(*b.base)[connectedDriver.ID] = connectedDriver

	mu.Unlock()

}

// удаляем водителя из базы подключенных водителей
func (b *ConnectedDriversBase) removeDriverFromBase(driverID uuid.UUID) {

	mu.Lock()

	// проверяем есть такой водитель в базе
	if driver, ok := (*b.base)[driverID]; ok {
		close(driver.DoneCH)             // посылаем, путем закрытия канала doneCH, сигнал о закрытии соединения
		close(driver.SendDataToDriverCH) // закрываем канал передачи данных водителю
		delete(*b.base, driverID)        // удаляем водителя из базы
	} else {
		b.logger.Info("Водителя нет в базе. driver_id: " + driverID.String())
	}

	mu.Unlock()

}
