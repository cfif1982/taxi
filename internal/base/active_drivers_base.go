package base

import (
	"encoding/json"

	"github.com/cfif1982/taxi/pkg/logger"
	"github.com/google/uuid"
)

// координаты водителя
type DriverGPS struct {
	DriverID  uuid.UUID `json:"driver_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}

// структура базы активных водителей
type ActiveDriversBase struct {
	base   *map[uuid.UUID]DriverGPS
	logger *logger.Logger
}

// Создаем базу
func CreateActiveDriversBase(logger *logger.Logger) (*ActiveDriversBase, error) {

	var base = make(map[uuid.UUID]DriverGPS)

	return &ActiveDriversBase{
		base:   &base,
		logger: logger,
	}, nil
}

// получить json строку всех активных водителей
func (b *ActiveDriversBase) GetAllActiveDriversString() (string, error) {

	driversString, err := json.Marshal(b.base)

	if err != nil {
		b.logger.Info("Неверный формат базы водителей:", err.Error())
	}

	return string(driversString), err
}

// обновить данные gps водителя в базе
func (b *ActiveDriversBase) UpdateGPSData(driverGPSData []byte) error {

	var driverGPS DriverGPS // храним полученные от водителя координаты

	err := json.Unmarshal(driverGPSData, &driverGPS)

	if err != nil {
		b.logger.Info("Неверный формат данных GPS от водителя:", err.Error())
		return err
	}

	// обращение к элементу map, переданной через указатель делается через (*map_name)
	(*b.base)[driverGPS.DriverID] = driverGPS

	return err
}
