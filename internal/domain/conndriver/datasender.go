package conndriver

type DataSenderToDriverInterface interface {
	SendDataToDriver(data []byte)
}
