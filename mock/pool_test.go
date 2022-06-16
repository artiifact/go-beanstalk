package mock_test

import (
	"context"
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPool_Open(t *testing.T) {
	p := &mock.Pool{}
	p.On("Open", context.Background()).Return(nil)

	require.NoError(t, p.Open(context.Background()))
}

func TestPool_Close(t *testing.T) {
	p := &mock.Pool{}
	p.On("Close", context.Background()).Return(nil)

	require.NoError(t, p.Close(context.Background()))
}

func TestPool_Get(t *testing.T) {
	p := &mock.Pool{}
	p.On("Get").Return(beanstalk.NewDefaultClient(mock.NewConn(nil, nil)), nil)

	c, err := p.Get()

	require.NotNil(t, c)
	require.Nil(t, err)
}

func TestPool_Put(t *testing.T) {
	c := beanstalk.NewDefaultClient(mock.NewConn(nil, nil))

	p := &mock.Pool{}
	p.On("Put", c).Return(nil)

	require.NoError(t, p.Put(c))
}

func TestPool_Len(t *testing.T) {
	p := &mock.Pool{}
	p.On("Len").Return(5)

	require.Equal(t, 5, p.Len())
}
