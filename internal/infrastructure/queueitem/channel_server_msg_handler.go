package queueitem

import (
	"errors"

	"github.com/cfif1982/taxi/internal/domain/queueitem"
)

// список возможных шибок
var (
	ErrServerChannelClosed = errors.New("server channel is closed")
)

// структура обработчика сообщений на сервере
// Реализована через каналы. Можно через RabbitMQ
type ChannelServerMsgHandler struct {
	dataCH chan *queueitem.QueueItem // канал, по которому будут передаваться данные на сервер от водителей
}

func NewChannelServerMsgHandler() *ChannelServerMsgHandler {

	return &ChannelServerMsgHandler{
		dataCH: make(chan *queueitem.QueueItem),
	}
}

// Получаем сообщения от водителя
func (c *ChannelServerMsgHandler) ReceiveMessageFromDriver() (*queueitem.QueueItem, error) {

	// ждем данные от водителей
	queueItem, ok := <-c.dataCH

	// Если канал данных закрыт, то ошибка
	if !ok {
		return nil, ErrServerChannelClosed
	}

	return queueItem, nil
}

// Посылаем данные на сервер
func (c *ChannelServerMsgHandler) SendMessageToServer(queueItem *queueitem.QueueItem) {

	// посылаем данные в канал
	c.dataCH <- queueItem
}
