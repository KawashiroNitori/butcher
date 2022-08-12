package butcher

type jobType int

const (
	jobTypeJob jobType = iota
	// jobTypeStop
)

type job[T any] struct {
	Type      jobType
	RetryTime int
	Payload   T
}
