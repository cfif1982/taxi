package conndriver

type DataReceiverFromDriverInterface interface {
	ReceiveDataFromDriver() (*ConnectedDriver, error)
}
