package conndriver

import (
	"sync"

	"github.com/cfif1982/taxi/internal/domain/conndriver"

	"github.com/google/uuid"
)

// репозиторий базы подключенных водителей
// QUESTION: в мапе хранить объекты или ссылки на объекты?
type InMemoryRepository struct {
	db map[uuid.UUID]*conndriver.ConnectedDriver // решил сделать map, а не слайс. В мапе быстрее будет искать по id нужного водителя
	mu sync.Mutex                                // мьютекс для синхронизации доступа к базе
}

func NewInMemoryRepo() *InMemoryRepository {
	return &InMemoryRepository{
		db: make(map[uuid.UUID]*conndriver.ConnectedDriver),
	}
}

// получить всех водителей из базы
func (r *InMemoryRepository) GetAllDrivers() (*[]*conndriver.ConnectedDriver, error) {

	// QUESTION: тут я правильно делаю? сначала блокирую, потом делаю копию, разблокирую и уже с копией работаю
	r.mu.Lock()

	// делаем копию базы, чтобы можно было по ней пробежаться не мешая добавлению данных
	baseCopy := r.db

	r.mu.Unlock()

	// слайс для отправки данных
	arrDrivers := []*conndriver.ConnectedDriver{}

	for _, v := range baseCopy {

		conDriver := conndriver.NewConnectedDriver(v.ID(), v.Latitude(), v.Longitude(), v.ReceivedDataTime())

		// сохраняем даные для формирования ответа сервера
		arrDrivers = append(arrDrivers, conDriver)
	}

	return &arrDrivers, nil
}

// обновить данные gps водителя в базе
// QUESTION: тоже самое: когда нужно возвращать ссылку на объект, а когда сам объект? Я путаюсь((( когда нужно передавать объект, а когда ссылку на него?
func (r *InMemoryRepository) UpdateDriver(connectedDriver *conndriver.ConnectedDriver) {

	r.mu.Lock()

	// памятка для меня: обращение к элементу map, переданной через указатель, делается через (*map_name)
	r.db[connectedDriver.ID()] = connectedDriver

	r.mu.Unlock()

}

// удаляем водителя из базы подключенных водителей
func (r *InMemoryRepository) RemoveDriver(driverID uuid.UUID) error {

	r.mu.Lock()

	// проверяем есть такой водитель в базе
	if _, ok := r.db[driverID]; ok {
		delete(r.db, driverID) // удаляем водителя из базы
	} else {
		return conndriver.ErrDriverNotFound
	}

	r.mu.Unlock()

	return nil
}
