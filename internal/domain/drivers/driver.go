package drivers

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// список возможных шибок
var (
	ErrTelephoneAlreadyExist = errors.New("telephone already exist")
	ErrWrongPassword         = errors.New("wrong login or password")
	ErrCookieError           = errors.New("cookie error")
	ErrDriverIsNotAuthorized = errors.New("driver is not authorized")
	ErrDriverIsNotFound      = errors.New("driver is not found")
	ErrInsufficientFunds     = errors.New("insufficient funds")          // недостаточно средств на балансе
	ErrHashGenerate          = errors.New("error while hash generation") // ошибка при генерации хэша пароля
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
func CreateDriver(routeID uuid.UUID, telephone, name, password string) (*Driver, error) {

	var zeroTime time.Time

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, ErrHashGenerate
	}

	return NewDriver(uuid.New(), routeID, telephone, name, hashedPassword, 0, zeroTime), nil
}

// Хэширование пароля
func hashPassword(password string) (string, error) {

	// хэшируем пароль
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// проверяем введенный пароль
func (d *Driver) CheckPassword(enteredPassword string) error {

	// сравниваем хэшированный пароль водителя и введенный пароль
	err := bcrypt.CompareHashAndPassword([]byte(d.password), []byte(enteredPassword))
	if err != nil {
		return ErrWrongPassword
	}

	return nil
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
