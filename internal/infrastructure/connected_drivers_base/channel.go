package connected_drivers_base

// подключенный к серверу водитель, т.е. водитель с которым активно соединение websocket
type ChannelSender struct {
	SendDataToDriverCH chan []byte   // канал, через который буду отсылаться данные водителю
	DoneCH             chan struct{} // канал, по которому будет передаваться сигнал о закрытии горутин приема и отправки данных
}
