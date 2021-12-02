package butcher

type Butcher interface {
	Run() error
}

type butcher struct {

}

func NewButcher() Butcher {
	return &butcher{}
}

func (b *butcher) Run() error {
	panic("implement me")
}

