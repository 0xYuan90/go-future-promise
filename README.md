## Future/Promise for Golang

[Future/Promise](http://dist-prog-book.com/chapter/2/futures.html)

## How to use

### Normal Case

```go
// func which will execute through future
fb := func() (interface{}, error) {
		time.Sleep(5 * time.Second)
		return 10, nil
}
// create future instance
f := future.New(fb)
// get result
v, err := f.Get()  

```

### With timeout

```go
// func which will execute through future
fb := func() (interface{}, error) {
		time.Sleep(5 * time.Second)
		return 10, nil
}
// create future instance
f := future.New(fb)
// get result with timeout setting
v, timeout, err := f.GetWithTimeout(3 * time.Second)

```

### Cancel before get result

```go
// func which will execute through future
f := future.New(func() (interface{}, error) {
		time.Sleep(10 * time.Second)
		return 5, nil
})
// cancel future after 2 seconds
go func() {
  time.Sleep(2 * time.Second)
  f.Cancel()
}()
// get result(after canceled)
result, err := f.Get()

```


### Then

```go
// use then to chain future together
f := future.New(func() (interface{}, error) {
	return 10, nil
}).Then(func(i interface{}) (interface{}, error) {
  return 2 * i.(int), nil
}).Then(func(i interface{}) (interface{}, error) {
  return 2 + i.(int), nil
})
result, err := f.Get()

```

