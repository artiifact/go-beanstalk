# Go Beanstalk

[![Build Status](https://github.com/IvanLutokhin/go-beanstalk/workflows/Test/badge.svg)](https://github.com/IvanLutokhin/go-beanstalk/actions?query=workflow%3ATest)
[![codecov](https://codecov.io/gh/IvanLutokhin/go-beanstalk/branch/master/graph/badge.svg)](https://codecov.io/gh/IvanLutokhin/go-beanstalk)
[![Go Reference](https://pkg.go.dev/badge/github.com/IvanLutokhin/go-beanstalk.svg)](https://pkg.go.dev/github.com/IvanLutokhin/go-beanstalk)

Go client for [beanstalkd](https://beanstalkd.github.io).

## Installation

```shell
go get github.com/IvanLutokhin/go-beanstalk
```

## Quick Start

### Producer:

```go
c, err := beanstalk.Dial("tcp", "127.0.0.1:11300")
if err != nil {
	panic(err)
}

id, err := c.Put(1, 0, 5*time.Second, []byte("example"))
if err != nil {
	panic(err)
}

fmt.Println(id) // output job id
```

### Consumer:
```go
c, err := beanstalk.Dial("tcp", "127.0.0.1:11300")
if err != nil {
	panic(err)
}

job, err := c.Reserve()
if err != nil {
	panic(err)
}

fmt.Println(job.ID) // output job id
fmt.Println(job.Data) // output job data
```

### Pool:
```go
p := beanstalk.NewPool(&beanstalk.PoolOptions{
	Dialer: func () (*Client, error) { return beanstalk.Dial("127.0.0.1:11300") },
	Logger: beanstalk.NopLogger,
	Capacity: 5,
	MaxAge: 0,
	IdleTimeout: 0,
})

// establish connections
if err = p.Open(); err != nil {
	panic(err)
}

// retrieve client
c, err := p.Get()
if err != nil {
	panic(err)
}

// use client
stats, err := c.Stats()
if err != nil {
	panic(err)
}

// return client
if err = p.Put(c); err != nil {
	panic(err)
}

// close connections
if err = p.Close(); err != nil {
	panic(err)
}
```

### HTTP Handler
```go
func GetServerStats() beanstalk.Handler {
	return beanstalk.HandlerFunc(func(c *beanstalk.Client, w http.ResponseWriter, r *http.Request) {
		stats, err := c.Stats()
		if err != nil {
			panic(err)
		}
		
		bytes, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	})
}

func main() {
    p := beanstalk.NewPool(&beanstalk.PoolOptions{
        Dialer: func () (*Client, error) { return beanstalk.Dial("127.0.0.1:11300") },
        Logger: beanstalk.NopLogger,
        Capacity: 5,
        MaxAge: 0,
        IdleTimeout: 0,
    })

    if err := p.Open(); err != nil {
        panic(err)
    }
	
    http.Handle("/stats", NewHTTPHandlerAdapter(p, GetServerStats()))
    http.ListenAndServe(":8090", nil)	
}
```

## License
[The MIT License (MIT)](LICENSE)
