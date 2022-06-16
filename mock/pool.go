package mock

import (
	"context"
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/stretchr/testify/mock"
)

type Pool struct {
	mock.Mock
}

func NewPool(client *Client) *Pool {
	pool := &Pool{}
	pool.On("Open", Anything).Return(nil)
	pool.On("Close", Anything).Return(nil)
	pool.On("Get").Return(client, nil)
	pool.On("Put", client).Return(nil)
	pool.On("Len", Anything).Return(1)

	return pool
}

func (p *Pool) Open(ctx context.Context) error {
	args := p.Called(ctx)

	return args.Error(0)
}

func (p *Pool) Close(ctx context.Context) error {
	args := p.Called(ctx)

	return args.Error(0)
}

func (p *Pool) Get() (beanstalk.Client, error) {
	args := p.Called()

	return args.Get(0).(beanstalk.Client), args.Error(1)
}

func (p *Pool) Put(client beanstalk.Client) error {
	args := p.Called(client)

	return args.Error(0)
}

func (p *Pool) Len() int {
	args := p.Called()

	return args.Int(0)
}
