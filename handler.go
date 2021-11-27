package beanstalk

import "net/http"

type Handler interface {
	ServeHTTP(client *Client, writer http.ResponseWriter, request *http.Request)
}

type HandlerFunc func(client *Client, writer http.ResponseWriter, request *http.Request)

func (f HandlerFunc) ServeHTTP(client *Client, writer http.ResponseWriter, request *http.Request) {
	f(client, writer, request)
}
