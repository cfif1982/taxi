package queueitem

import (
	"errors"
)

// список возможных шибок
var (
	ErrDriverChannelClosed = errors.New("driver channel is closed")
)

// структура обработчика сообщений у водителя
// Реализована через каналы. Можно через RabbitMQ
type ChannelDriverMsgHandler struct {
	dataCH chan []byte   // канал, через который буду отсылаться данные водителю с сервера
	doneCH chan struct{} // канал, по которому с сервера будет передаваться сигнал о закрытии хэндлера
}

func NewChannelDriverMsgHandler() *ChannelDriverMsgHandler {

	return &ChannelDriverMsgHandler{
		dataCH: make(chan []byte),
		doneCH: make(chan struct{}),
	}
}

// Получаем сообщения от сервера
func (c *ChannelDriverMsgHandler) ReceiveMessageFromServer() ([]byte, error) {

	// Если хэндлер будет закрыт, то завершаем работу
	for {
		select {
		case <-c.doneCH: // закрытие канала done
			return nil, ErrDriverChannelClosed

		case receivedData, ok := <-c.dataCH: // ждем данные от сервера

			// Если канал данных закрыт, то ошибка
			if !ok {
				return nil, ErrDriverChannelClosed
			}

			return receivedData, nil
		}
	}
}

// посылаем данные водителю
func (c *ChannelDriverMsgHandler) SendMessageToDriver(data []byte) {

	// посылаем данные в канал
	c.dataCH <- data
}

// Ждем сигнал на закрытие хэндлера
func (c *ChannelDriverMsgHandler) WaitCloseSignal() {

	// если сигнальный канал закрыт, то закрываем хэндлер
	_, _ = <-c.doneCH

	close(c.dataCH)
}

// закрываем хэндлер
func (c *ChannelDriverMsgHandler) CloseHandler() {

	// для закрытия хэндлера, нужно закрыть сигнальный канал
	close(c.doneCH)
}
