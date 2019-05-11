package sender

type Sender interface {
	Run() error
}

type QpsSender struct {
	rate     int
	taskId   int
	taskArgs interface{}
}
