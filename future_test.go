package future_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	future "github.com/fatfish90/go-future-promise"
	"github.com/stretchr/testify/assert"
)

func TestFutureGet(t *testing.T) {
	fb := func() (interface{}, error) {
		time.Sleep(5 * time.Second)
		return 10, nil
	}
	f := future.New(fb)
	start := time.Now()
	v, err := f.Get()
	end := time.Now()
	dur := end.Unix() - start.Unix()
	fmt.Println(v, err, dur)
	assert.Nil(t, err)
	assert.Equal(t, 10, v)
	assert.Equal(t, int64(5), dur)
}

func TestFutureGetWithTimeout(t *testing.T) {
	fb := func() (interface{}, error) {
		time.Sleep(5 * time.Second)
		return 10, nil
	}
	f := future.New(fb)
	start := time.Now()
	v, timeout, err := f.GetWithTimeout(3 * time.Second)
	end := time.Now()
	dur := end.Unix() - start.Unix()
	fmt.Println(v, timeout, err, dur)
	assert.Nil(t, err)
	assert.Nil(t, v)
	assert.True(t, timeout)
	assert.True(t, dur < 5)

	start2 := time.Now()
	v, timeout, err = f.GetWithTimeout(10 * time.Second)
	end2 := time.Now()
	dur2 := end2.Unix() - start2.Unix()
	fmt.Println(v, timeout, err, dur2)
	assert.Nil(t, err)
	assert.Equal(t, 10, v)
	assert.False(t, timeout)
	assert.True(t, dur < int64(10))
}

func TestThen(t *testing.T) {
	f := future.New(func() (interface{}, error) {
		return 10, nil
	}).Then(func(i interface{}) (interface{}, error) {
		return 2 * i.(int), nil
	}).Then(func(i interface{}) (interface{}, error) {
		return 2 + i.(int), nil
	})
	result, err := f.Get()
	fmt.Println(result, err)
	assert.Nil(t, err)
	assert.Equal(t, 22, result)

	g := future.New(func() (interface{}, error) {
		return nil, errors.New("This is an error")
	}).Then(func(i interface{}) (interface{}, error) {
		return 2 * i.(int), nil
	})
	result, err = g.Get()
	fmt.Println(result, err)
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "This is an error", err.Error())

	h := future.New(func() (interface{}, error) {
		return 10, nil
	}).Then(func(i interface{}) (interface{}, error) {
		return nil, errors.New("This is also an error")
	})
	result, err = h.Get()
	fmt.Println(result, err)
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "This is also an error", err.Error())
}

func TestCancel(t *testing.T) {
	f := future.New(func() (interface{}, error) {
		time.Sleep(10 * time.Second)
		return 5, nil
	})
	go func() {
		time.Sleep(2 * time.Second)
		f.Cancel()
	}()
	result, err := f.Get()
	fmt.Println(result, err, f.IsCancelled())
	assert.Nil(t, result)
	assert.Nil(t, err)
	assert.True(t, f.IsCancelled())
}
