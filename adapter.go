package beanstalk

import "net/http"

type HTTPHandlerAdapter struct {
	pool    *Pool
	handler Handler
}

func NewHTTPHandlerAdapter(pool *Pool, handler Handler) *HTTPHandlerAdapter {
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
