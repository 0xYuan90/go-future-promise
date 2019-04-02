package future

import (
	"context"
	"time"
)

type Future interface {
	Cancel()

	IsCancelled() bool

	Get() (interface{}, error)

	GetWithTimeout(d time.Duration) (interface{}, bool, error)

	Then(func(interface{}) (interface{}, error)) Future
}

func New(inFunc func() (interface{}, error)) Future {
	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	return execute(cancelCtx.Done(), cancelFunc, inFunc)
}

type FutureImpl struct {
	done       chan struct{}
	cancelChan <-chan struct{}
	cancelFunc context.CancelFunc
	val        interface{}
	err        error
}

func (f *FutureImpl) Cancel() {
	select {
	case <-f.done:
	case <-f.cancelChan:
		return
	default:
		f.cancelFunc()
	}
}

func (f *FutureImpl) IsCancelled() bool {
	select {
	case <-f.cancelChan:
		return true
	default:
		return false
	}
}

func (f *FutureImpl) Get() (interface{}, error) {
	select {
	case <-f.done:
		return f.val, f.err
	case <-f.cancelChan:
		return nil, nil
	}
}

func (f *FutureImpl) GetWithTimeout(d time.Duration) (interface{}, bool, error) {
	select {
	case <-f.done:
		val, err := f.Get()
		return val, false, err
	case <-time.After(d):
		return nil, true, nil
	default:
		return nil, false, nil
	}
}

func (f *FutureImpl) Then(next func(interface{}) (interface{}, error)) Future {
	nextFuture := execute(f.cancelChan, f.cancelFunc, func() (interface{}, error) {
		result, err := f.Get()
		if f.IsCancelled() || err != nil {
			return result, err
		}
		return next(result)
	})
	return nextFuture
}

func execute(cancelChan <-chan struct{}, cancelFunc context.CancelFunc, inFunc func() (interface{}, error)) Future {
	f := FutureImpl{
		done:       make(chan struct{}),
		cancelChan: cancelChan,
		cancelFunc: cancelFunc,
	}
	go func() {
		go func() {
			f.val, f.err = inFunc()
			close(f.done)
		}()
		select {
		case <-f.done:
			//do nothing, just waiting to see which will happen first
		case <-f.cancelChan:
			//do nothing, leave val and err nil
		}
	}()
	return &f
}
