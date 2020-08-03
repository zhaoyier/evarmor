package network

type WorkerPool struct {
	workers   []*worker
	closeChan chan struct{}
}

// WorkerPoolInstance returns the global pool.
func WorkerPoolInstance() *WorkerPool {
	return globalWorkerPool
}

func newWorkerPool(vol int) *WorkerPool {
	if vol <= 0 {
		vol = defaultWorkersNum
	}

	pool := &WorkerPool{
		workers:   make([]*worker, vol),
		closeChan: make(chan struct{}),
	}

	for i := range pool.workers {
		pool.workers[i] = newWorker(i, 1024, pool.closeChan)
		if pool.workers[i] == nil {
			panic("worker nil")
		}
	}
	return pool
}

func newWorker(i int, c int, closeChan chan struct{}) *worker {
	w := &worker{
		index:        i,
		callbackChan: make(chan workerFunc, c),
		closeChan:    closeChan,
	}
	go w.start()
	return w
}

type worker struct {
	index        int
	callbackChan chan workerFunc
	closeChan    chan struct{}
}

func (w *worker) start() {
	for {
		select {
		case <-w.closeChan:
			return
		case cb := <-w.callbackChan:
			cb()
		}
	}
}

func (w *worker) put(cb workerFunc) error {
	select {
	case w.callbackChan <- cb:
		return nil
	default:
		return ErrWouldBlock
	}
}
