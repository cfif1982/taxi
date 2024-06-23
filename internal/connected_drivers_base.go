package internal

import (
	"time"

	"github.com/cfif1982/taxi/internal/domain/conndriver"
	"github.com/cfif1982/taxi/pkg/logger"
)

const SendDataPeriod = time.Second * 2 // частота отсылки данных водителю

// база подключенных к серверу водителей
type ConnectedDriversBase struct {
	dataReceiver *conndriver.ReceiverDataFromDriver
	logger       *logger.Logger
}

// Конструктор ConnectedDriversBase
// Question: вот тут тоже непонятно - возвращатьт указатель или саму стурктуру? При создании сервера - возвращаем саму стурктуру.
// а при создании водителя - возвращаем указатель на водителя
func NewConnectedDriversBase(dataReceiver *conndriver.ReceiverDataFromDriver, logger *logger.Logger) ConnectedDriversBase {
	return ConnectedDriversBase{
		dataReceiver: dataReceiver,
		logger:       logger,
	}
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

		// обновляем данные водителя в базе. Если его там нет, то добавляем
		case connectedDriver := <-b.ReceiveDataFromDriverCH:
			b.updateDriversData(connectedDriver)

			// удаляем неактивного водителя из базы. Пока не знаю - нужно это или нет
			// водитель хочет завершить работу. Жмет у себя кнопку закрыть приложение (ну или отключиться) - вот в этом слугае от него и идет такой сигнал.
			// case driverID := <-b.RemoveDriverCH:
			// 	b.removeDriverFromBase(driverID)
		}
	}
}
