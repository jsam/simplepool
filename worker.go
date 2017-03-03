package simplepool

type Worker struct {
	pool    chan *Worker
	channel chan Job
	stop    chan bool
}

func (w *Worker) start() {
	go func() {
		var job Job
		for {
			w.pool <- w
			select {
			case job = <-w.channel:
				job()
			case stop := <-w.stop:
				if stop {
					w.stop <- true
					return
				}
			}
		}
	}()
}

func newWorker(p chan *Worker) *Worker {
	return &Worker{
		pool:    p,
		channel: make(chan Job),
		stop:    make(chan bool),
	}
}
