package beanstalk

import (
	"net/http"
)

// Handler

type Handler interface {
	ServeHTTP(client Client, writer http.ResponseWriter, request *http.Request)
}

type HandlerFunc func(client Client, writer http.ResponseWriter, request *http.Request)

func (f HandlerFunc) ServeHTTP(client Client, writer http.ResponseWriter, request *http.Request) {
	f(client, writer, request)
}

// Adapter

type HTTPHandlerAdapter struct {
	pool    Pool
	handler Handler
}

func NewHTTPHandlerAdapter(pool Pool, handler Handler) *HTTPHandlerAdapter {
	return &HTTPHandlerAdapter{
		pool:    pool,
		handler: handler,
	}
}

func (a *HTTPHandlerAdapter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	client, err := a.pool.Get()
	if err != nil {
		panic(err)
	}

	defer a.pool.Put(client)

	a.handler.ServeHTTP(client, writer, request)
}
