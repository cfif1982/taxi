package conndriver

import (
	"fmt"

	"github.com/cfif1982/taxi/internal/domain/conndriver"
)

// канал для получения данных от подключенного водителя
type ChannelReceiver struct {
	ReceiveDataFromDriverCH chan *conndriver.ConnectedDriver // канал, по которому будут передаваться данные от водителейых
}

func (c *ChannelReceiver) NewChannelReceiver() *ChannelReceiver {

	return &ChannelReceiver{
		ReceiveDataFromDriverCH: make(chan *conndriver.ConnectedDriver),
	}

}

func (c *ChannelReceiver) ReceiveDataFromDriver() (*conndriver.ConnectedDriver, error) {

	connDriver, ok := <-c.ReceiveDataFromDriverCH

	if !ok {
		return nil, fmt.Errorf("channel closed")
	}

	return connDriver, nil
}
