package drivers

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// список возможных шибок
var (
	ErrTelephoneAlreadyExist = errors.New("telephone already exist")
	ErrWrongPassword         = errors.New("wrong login or password")
	ErrCookieError           = errors.New("cookie error")
	ErrDriverIsNotAuthorized = errors.New("driver is not authorized")
	ErrDriverIsNotFound      = errors.New("driver is not found")
	ErrInsufficientFunds     = errors.New("insufficient funds ") // недостаточно средств на балансе
)

// структура для хранения водителя
type Driver struct {
	id           uuid.UUID
	routeID      uuid.UUID
	telephone    string
	password     string
	name         string
	balance      int
	lastPaidDate time.Time
}

// создаем новый объект
// нужна для использвания в других пакетах
func NewDriver(id, routeID uuid.UUID, telephone, name, password string, balance int, lastPaidDate time.Time) *Driver {
	return &Driver{
		id:           id,
		routeID:      routeID,
		telephone:    telephone,
		name:         name,
		password:     password,
		balance:      balance,
		lastPaidDate: lastPaidDate,
	}
}

// Создаем водителя
func CreateDriver(routeID uuid.UUID, telephone, name, password string) *Driver {

	var zeroTime time.Time

	return NewDriver(uuid.New(), routeID, telephone, name, password, 0, zeroTime)
}

// увеличить баланс
func (d *Driver) IncreaseBalance(summa int) error {

	d.balance += summa

	return nil
}

// Уменьшить баланс
func (d *Driver) ReduceBalance(summa int) error {

	// Если недостаточно средств - ошибка
	if d.balance < summa {
		return ErrInsufficientFunds
	}

	d.balance -= summa

	// изменяем дату последней оплаты на сегодняшнюю дату
	currentDate := time.Now()

	d.lastPaidDate = currentDate

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

// возвращщаем поле lastPaidDate
func (d *Driver) LastPaidDate() time.Time {
	return d.lastPaidDate
}
