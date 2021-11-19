# Go Beanstalk

Go client for [beanstalkd](https://beanstalkd.github.io).

## Installation

```shell
go get github.com/IvanLutokhin/go-beanstalk
```

## Quick Start

Producer:

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

Consumer:
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

Pool:
```go
p, err := beanstalk.NewDefaultPool("127.0.0.1:11300", 5, 10)
if err != nil {
	panic(err)
}

c, err := p.Get()
if err != nil {
	panic(err)
}

// use client

if err = p.Put(c); err != nil {
	panic(err)
}
```
## License
[The MIT License (MIT)](LICENSE)
