package simplepool

import "sync"

// Job is type of atomic object which gets passed around.
type Job func(args ...interface{})

// Pool is type for manipulating and dispatching jobs.
type Pool struct {
	Jobs       chan Job
	dispatcher *Dispatcher
	wg         sync.WaitGroup
}

// NewPool creates new Pool object with specific number of worker and queue length.
func NewPool(workersCount int, queueSize int) *Pool {
	queueChannel := make(chan Job, queueSize)
	workersChannel := make(chan *Worker, workersCount)

	pool := &Pool{
		Jobs:       queueChannel,
		dispatcher: newDispatcher(workersChannel, queueChannel),
	}

	return pool
}

// JobDone needs to be called when job finishes it's work.
func (p *Pool) JobDone() {
	p.wg.Done()
}

// WaitCount is configuration option which enables us to configure how many jobs
// we should wait when calling WaitAll.
func (p *Pool) WaitCount(count int) {
	p.wg.Add(count)
}

// WaitAll is sync method for all workers.
func (p *Pool) WaitAll() {
	p.wg.Wait()
}

// Release will free up all resources used by the pool.
func (p *Pool) Release() {
	p.dispatcher.stop <- true
	<-p.dispatcher.stop
}

// Enqueue will enqeueu new arbitrary function to the pool and execute it as a job.
func (p *Pool) Enqueue(job Job, args ...interface{}) {
	arg := args
	p.Jobs <- func(args ...interface{}) {
		defer p.JobDone()
		job(arg...)
	}
}
