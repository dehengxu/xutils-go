package pkg

type CommonResultInterface interface {
	GetData() interface{}
	GetError() error
}

type ResultBox[T any] struct {
	CommonResultInterface
	data interface{}
	err  error
}

func (r *ResultBox[T]) Value() T {
	return r.data.(T)
}

func (r *ResultBox[T]) ValueOrDefault(d T) T {
	if r.err != nil || r.data == nil {
		return d
	}
	return r.data.(T)
}

func (r *ResultBox[T]) Error() error {
	return r.err
}

func Wrap[T any](input T, err error) *ResultBox[T] {
	return &ResultBox[T]{data: input, err: err}
}
