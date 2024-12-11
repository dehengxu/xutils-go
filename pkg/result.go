package pkg

import "errors"

type CommonResultInterface[T any] interface {
	Value() T
	Error() error
	OnFailed(handler func(e error)) CommonResultInterface[T]
	OnSuccess(handler func(data T)) CommonResultInterface[T]
	OnFinal(handler func())
}

type ResultBox[T any] struct {
	CommonResultInterface[T]
	data T
	err  error
	isok bool
}

func (r *ResultBox[T]) ValueOrDefault(d T) T {
	if r.err != nil {
		return d
	}
	return r.data
}

func (r *ResultBox[T]) HasValue() bool {
	return r.isok
}

func (r *ResultBox[T]) Value() T {
	if r.HasValue() {
		return r.data
	} else {
		var t T
		return t
	}
}

func (r *ResultBox[T]) Error() error {
	return r.err
}

func (r *ResultBox[T]) OnFailed(handler func(e error)) CommonResultInterface[T] {
	if handler == nil {
		return r
	}
	if !r.HasValue() {
		handler(r.err)
	}
	return r
}

func (r *ResultBox[T]) OnSuccess(handler func(data T)) CommonResultInterface[T] {
	if handler == nil {
		return r
	}
	if r.HasValue() {
		handler(r.data)
	}
	return r
}

func (r *ResultBox[T]) OnFinal(handler func()) {
	if handler != nil {
		handler()
	}
}

func Wrap[T any](input T, err interface{}) *ResultBox[T] {
	if value, ok := err.(error); ok {
		return &ResultBox[T]{
			data: input,
			err:  value,
			isok: func(e error) bool { return e == nil }(value),
		}
	} else if value1, ok := err.(bool); ok {
		return &ResultBox[T]{
			data: input,
			err: func(b bool) error {
				if b {
					return nil
				} else {
					return errors.New("error :false")
				}
			}(value1),
			isok: value1}
	}
	return &ResultBox[T]{data: input, err: nil, isok: true}
}
