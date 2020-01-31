package worker

type Worker struct {
	pool Pool
}

func NewWorker() *Worker {
	return &Worker{pool: NewPool()}
}

//func (w *Worker) Call(call Call, options ...CallOption) error {
//	runner := w.pool.Get()
//	defer w.pool.Put(runner)
//	return runner.Call(call, options...)
//}
