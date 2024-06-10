package drivers

// "github.com/google/uuid"

// список возможных шибок
// var (
// 	ErrLinkNotFound    = errors.New("link not found")
// 	ErrKeyAlreadyExist = errors.New("key already exist")
// 	ErrURLAlreadyExist = errors.New("url already exist")
// )

// структура для хранения водителя
type Driver struct {
	id        int
	fio       string
	latitude  uint64
	longitude uint64
	balance   int
}

// создаем новый объект
// нужна для использвания в других пакетах
func NewDriver(id int, fio string, latitude, longitude uint64, balance int) (*Driver, error) {
	return &Driver{
		id:        id,
		fio:       fio,
		latitude:  latitude,
		longitude: longitude,
		balance:   balance,
	}, nil
}

// Создаем новую ССЫЛКУ
func CreateDriver(fio string) (*Driver, error) {

	// QUESTION: нужно ли здесь генерировать ID? Я хочу чтобы этот id назначала сама БД, чтобы не проводить проверку на существование такого id,
	// т.к. этот id должен быть уникальным
	return NewDriver(0, fio, 0, 0, 0)
}

// // генерируем key
// func generateKey() string {

// 	// генерируем случайный код типа string
// 	uuid := uuid.NewString()[:8]

// 	return uuid
// }

// // возвращщаем поле key
// func (l *Link) Key() string {
// 	return l.key
// }

// // возвращщаем поле URL
// func (l *Link) URL() string {
// 	return l.url
// }

// // возвращщаем поле UserID
// func (l *Link) UserID() int {
// 	return l.userID
// }

// // возвращщаем поле deletedFlag
// func (l *Link) DeletedFlag() bool {
// 	return l.deletedFlag
// }
