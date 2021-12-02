package butcher

type Executor interface {
	GenerateJob(chan<- interface{}) error
	Task(job interface{}) error
}
