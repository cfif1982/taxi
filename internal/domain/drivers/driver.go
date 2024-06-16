package drivers

import (
	"errors"

	"github.com/google/uuid"
)

// список возможных шибок
var (
	ErrTelephoneAlreadyExist = errors.New("telephone already exist")
	ErrWrongPassword         = errors.New("wrong login or password")
)

// структура для хранения водителя
type Driver struct {
	id        uuid.UUID
	routeID   uuid.UUID
	telephone string
	password  string
	name      string
	latitude  uint64
	longitude uint64
	balance   int
}

// создаем новый объект
// нужна для использвания в других пакетах
func NewDriver(id, routeID uuid.UUID, telephone, name, password string, balance int) (*Driver, error) {
	return &Driver{
		id:        id,
		routeID:   routeID,
		telephone: telephone,
		name:      name,
		password:  password,
		balance:   balance,
	}, nil
}

// Создаем водителя
func CreateDriver(routeID uuid.UUID, telephone, name, password string) (*Driver, error) {

	return NewDriver(uuid.New(), routeID, telephone, name, password, 0)
}

// возвращщаем поле ID
func (d *Driver) IncreaseBalance(summa int) error {

	d.balance += summa

	return nil
}

// возвращщаем поле ID
func (d *Driver) ID() uuid.UUID {
	return d.id
}

// возвращщаем поле routeID
func (d *Driver) RouteID() uuid.UUID {
	return d.routeID
}

// возвращщаем поле telephone
func (d *Driver) Telephone() string {
	return d.telephone
}

// возвращщаем поле name
func (d *Driver) Name() string {
	return d.name
}

// возвращщаем поле password
func (d *Driver) Password() string {
	return d.password
}

// возвращщаем поле balance
func (d *Driver) Balance() int {
	return d.balance
}
