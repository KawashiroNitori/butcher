package butcher

type Executor interface {
	GenerateJob(chan<- interface{}) error
	Task(workerID int, job interface{}) error
}
