package pkg

import "fmt"

type CommonResultInterface interface {
	GetData() interface{}
	GetError() error
}

type ResultBox[T any] struct {
	CommonResultInterface
	data interface{}
	err  error
	isok bool
}

func (r *ResultBox[T]) OnSuccess(f func(T)) *ResultBox[T] {
	if r.err == nil || r.isok {
		f(r.data.(T))
	}
	return r
}

func (r *ResultBox[T]) OnFailure(f func(error)) *ResultBox[T] {
	if r.err != nil {
		f(r.err)
	}
	if !r.isok {
		f(fmt.Errorf("result is false"))
	}
	return r
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

func (r *ResultBox[T]) IsOK() bool {
	return r.err == nil && r.isok
}

func (r *ResultBox[T]) Error() error {
	if r.err != nil {
		return r.err
	}
	if !r.isok {
		return fmt.Errorf("result is false")
	}
	return nil
}

func Wrap[T any](input T, err error) *ResultBox[T] {
	return &ResultBox[T]{data: input, err: err}
}
