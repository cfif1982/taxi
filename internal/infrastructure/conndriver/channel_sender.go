package conndriver

// "github.com/cfif1982/taxi/internal/domain/conndriver"

// канал для отсылки данных подключенному водителю
type ChannelSender struct {
	SendDataToDriverCH chan []byte   // канал, через который буду отсылаться данные водителю
	DoneCH             chan struct{} // канал, по которому будет передаваться сигнал о закрытии горутин приема и отправки данных
}

func (c *ChannelSender) SendDataToDriver(data []byte) {
	c.SendDataToDriverCH <- data
}
