package workers

import (
	"fmt"
)

type Worker struct {
	id         int
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func InitJobQueue() {
	JobQueue = make(chan Job)
}

func NewWorker(workerPool chan chan Job, id int) Worker {
	return Worker{
		id:         id,
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				fmt.Println("Worker id is ", w.id)
				if err := job.DoJob(); err != nil {
					//log.Errorf("Error when join job: %s", err.Error())
					fmt.Printf("Error when join job: %s\n", err.Error())
				}
			case <-w.quit:
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
