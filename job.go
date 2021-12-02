package butcher

type jobType int

const (
	jobTypeJob jobType = iota
	jobTypeStop
)

type job struct {
	Type      jobType
	RetryTime int
	Payload   interface{}
}
