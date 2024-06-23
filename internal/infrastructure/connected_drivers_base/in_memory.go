package connected_drivers_base

import (
	"sync"

	base "github.com/cfif1982/taxi/internal/domain/connected_drivers_base"

	"github.com/google/uuid"
)

// репозиторий базы подключенных водителей
type InMemoryRepository struct {
	base map[uuid.UUID]base.ConnectedDriver // решил сделать map, а не слайс. В мапе быстрее будет искать по id нужного водителя
	mu   sync.Mutex                         // мьютекс для синхронизации доступа к базе
}

func NewInMemoryRepo() *InMemoryRepository {
	return &InMemoryRepository{
		base: make(map[uuid.UUID]base.ConnectedDriver),
	}
}

// получить всех водителей из базы
// QUESTION: когда нужно возвращать ссылку на объект, а когда сам объект? Я путаюсь((( когда нужно передавать объект, а когда ссылку на него?
func (r *InMemoryRepository) GetAllDrivers() (*[]base.ConnectedDriver, error) {

	// QUESTION: тут я правильно делаю? сначала блокирую, потом делаю копию, разблокирую и уже с копией работаю
	r.mu.Lock()

	// делаем копию базы, чтобы можно было по ней пробежаться не мешая добавлению данных
	baseCopy := r.base

	r.mu.Unlock()

	// слайс для отправки данных
	arrDrivers := []base.ConnectedDriver{}

	for _, v := range baseCopy {

		conDriver := base.NewConnectedDriver(v.ID(), v.Latitude(), v.Longitude(), v.ReceivedDataTime())

		// сохраняем даные для формирования ответа сервера
		arrDrivers = append(arrDrivers, *conDriver)
	}

	return &arrDrivers, nil
}

// обновить данные gps водителя в базе
// QUESTION: тоже самое: когда нужно возвращать ссылку на объект, а когда сам объект? Я путаюсь((( когда нужно передавать объект, а когда ссылку на него?
func (r *InMemoryRepository) UpdateDriver(connectedDriver *base.ConnectedDriver) {

	r.mu.Lock()

	// памятка для меня: обращение к элементу map, переданной через указатель, делается через (*map_name)
	r.base[connectedDriver.ID()] = *connectedDriver

	r.mu.Unlock()

}

// удаляем водителя из базы подключенных водителей
func (r *InMemoryRepository) RemoveDriver(driverID uuid.UUID) error {

	r.mu.Lock()

	// проверяем есть такой водитель в базе
	if _, ok := r.base[driverID]; ok {
		delete(r.base, driverID) // удаляем водителя из базы
	} else {
		return base.ErrDriverNotFound
	}

	r.mu.Unlock()

	return nil
}
