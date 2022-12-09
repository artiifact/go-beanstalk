package beanstalk_test

import (
	"context"
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewDefaultPool(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		pool := beanstalk.NewDefaultPool(&beanstalk.PoolOptions{
			Dialer: func() (*beanstalk.Client, error) {
				return beanstalk.NewClient(mock.NewConn(nil, nil)), nil
			},
		})

		require.Equal(t, 0, pool.Len())

		// opens pool
		require.NoError(t, pool.Open(context.Background()))

		require.Equal(t, 1, pool.Len())

		// gets client from pool
		client, err := pool.Get()

		require.Nil(t, err)
		require.NotNil(t, client)

		// puts client in pool
		require.NoError(t, pool.Put(client))

		// closes pool
		require.NoError(t, pool.Close(context.Background()))

		require.Equal(t, 0, pool.Len())
	})

	t.Run("dialer not specified", func(t *testing.T) {
		pool := beanstalk.NewDefaultPool(&beanstalk.PoolOptions{})

		require.Equal(t, 0, pool.Len())

		// opens pool
		require.NoError(t, pool.Open(context.Background()))

		require.Equal(t, 0, pool.Len())

		// gets client from pool
		client, err := pool.Get()

		require.Equal(t, beanstalk.ErrDialerNotSpecified, err)
		require.Nil(t, client)

		// closes pool
		require.NoError(t, pool.Close(context.Background()))

		require.Equal(t, 0, pool.Len())
	})
}

func TestDefaultPool_Open(t *testing.T) {
	t.Run("timeout", func(t *testing.T) {
		pool := beanstalk.NewDefaultPool(&beanstalk.PoolOptions{
			Dialer: func() (*beanstalk.Client, error) {
				return beanstalk.NewClient(mock.NewConn(nil, nil)), nil
			},
			Logger:      beanstalk.NopLogger,
			Capacity:    3,
			MaxAge:      0,
			IdleTimeout: 0,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 0)

		defer cancel()

		require.Error(t, pool.Open(ctx))
	})

	t.Run("already opened", func(t *testing.T) {
		pool := beanstalk.NewDefaultPool(&beanstalk.PoolOptions{
			Dialer: func() (*beanstalk.Client, error) {
				return beanstalk.NewClient(mock.NewConn(nil, nil)), nil
			},
			Logger:      beanstalk.NopLogger,
			Capacity:    3,
			MaxAge:      0,
			IdleTimeout: 0,
		})

		require.Equal(t, 0, pool.Len())

		// opens pool
		require.NoError(t, pool.Open(context.Background()))

		require.Equal(t, 3, pool.Len())

		// retries to open pool
		require.Error(t, pool.Open(context.Background()))

		require.Equal(t, 3, pool.Len())
	})
}

func TestDefaultPool_Close(t *testing.T) {
	t.Run("timeout", func(t *testing.T) {
		pool := beanstalk.NewDefaultPool(&beanstalk.PoolOptions{
			Dialer: func() (*beanstalk.Client, error) {
				return beanstalk.NewClient(mock.NewConn(nil, nil)), nil
			},
			Logger:      beanstalk.NopLogger,
			Capacity:    5,
			MaxAge:      0,
			IdleTimeout: 0,
		})

		require.Equal(t, 0, pool.Len())

		// opens pool
		require.NoError(t, pool.Open(context.Background()))

		require.Equal(t, 5, pool.Len())

		// closes pool
		ctx, cancel := context.WithTimeout(context.Background(), 0)

		defer cancel()

		require.Error(t, pool.Close(ctx))
	})

	t.Run("closed", func(t *testing.T) {
		pool := beanstalk.NewDefaultPool(&beanstalk.PoolOptions{
			Dialer: func() (*beanstalk.Client, error) {
				return beanstalk.NewClient(mock.NewConn(nil, nil)), nil
			},
			Logger:      beanstalk.NopLogger,
			Capacity:    5,
			MaxAge:      0,
			IdleTimeout: 0,
		})

		require.Equal(t, 0, pool.Len())

		// opens pool
		require.NoError(t, pool.Open(context.Background()))

		require.Equal(t, 5, pool.Len())

		// closes pool
		require.NoError(t, pool.Close(context.Background()))

		require.Equal(t, 0, pool.Len())

		// retries to close pool
		require.Error(t, pool.Close(context.Background()))
	})
}

func TestDefaultPool_Get(t *testing.T) {
	t.Run("max age", func(t *testing.T) {
		pool := beanstalk.NewDefaultPool(&beanstalk.PoolOptions{
			Dialer: func() (*beanstalk.Client, error) {
				return beanstalk.NewClient(mock.NewConn(nil, nil)), nil
			},
			Logger:      beanstalk.NopLogger,
			Capacity:    1,
			MaxAge:      1 * time.Second,
			IdleTimeout: 0,
		})

		require.Equal(t, 0, pool.Len())

		// opens pool
		require.NoError(t, pool.Open(context.Background()))

		require.Equal(t, 1, pool.Len())

		// waiting timeout
		time.Sleep(1 * time.Second)

		// gets client by factory
		client, err := pool.Get()

		require.Nil(t, err)
		require.NotNil(t, client)

		// closes pool
		require.NoError(t, pool.Close(context.Background()))

		require.Equal(t, 0, pool.Len())
	})

	t.Run("idle timeout", func(t *testing.T) {
		pool := beanstalk.NewDefaultPool(&beanstalk.PoolOptions{
			Dialer: func() (*beanstalk.Client, error) {
				return beanstalk.NewClient(mock.NewConn(nil, nil)), nil
			},
			Logger:      beanstalk.NopLogger,
			Capacity:    1,
			MaxAge:      0,
			IdleTimeout: 1 * time.Second,
		})

		require.Equal(t, 0, pool.Len())

		// opens pool
		require.NoError(t, pool.Open(context.Background()))

		require.Equal(t, 1, pool.Len())

		// waiting timeout
		time.Sleep(1 * time.Second)

		// gets client by factory
		client, err := pool.Get()

		require.Nil(t, err)
		require.NotNil(t, client)

		// closes pool
		require.NoError(t, pool.Close(context.Background()))

		require.Equal(t, 0, pool.Len())
	})

	t.Run("closed", func(t *testing.T) {
		pool := beanstalk.NewDefaultPool(&beanstalk.PoolOptions{
			Dialer: func() (*beanstalk.Client, error) {
				client := beanstalk.NewClient(mock.NewConn(nil, nil))
				if err := client.Close(); err != nil {
					t.Fatal(err)
				}

				return client, nil
			},
			Logger:      beanstalk.NopLogger,
			Capacity:    1,
			MaxAge:      0,
			IdleTimeout: 0,
		})

		require.Equal(t, 0, pool.Len())

		// opens pool
		require.NoError(t, pool.Open(context.Background()))

		require.Equal(t, 1, pool.Len())

		// gets client by factory
		client, err := pool.Get()

		require.Nil(t, err)
		require.NotNil(t, client)

		// closes pool
		require.NoError(t, pool.Close(context.Background()))

		require.Equal(t, 0, pool.Len())
	})
}

func TestDefaultPool_Put(t *testing.T) {
	t.Run("closed pool", func(t *testing.T) {
		pool := beanstalk.NewDefaultPool(&beanstalk.PoolOptions{
			Dialer: func() (*beanstalk.Client, error) {
				return beanstalk.NewClient(mock.NewConn(nil, nil)), nil
			},
			Logger:      beanstalk.NopLogger,
			Capacity:    1,
			MaxAge:      0,
			IdleTimeout: 0,
		})

		require.Equal(t, 0, pool.Len())

		// opens pool
		require.NoError(t, pool.Open(context.Background()))

		require.Equal(t, 1, pool.Len())

		// gets client
		client, err := pool.Get()

		require.Nil(t, err)
		require.NotNil(t, client)

		// closes pool
		require.NoError(t, pool.Close(context.Background()))

		require.Equal(t, 0, pool.Len())

		// puts client in closed pool
		require.Error(t, pool.Put(client))
	})

	t.Run("closed client", func(t *testing.T) {
		pool := beanstalk.NewDefaultPool(&beanstalk.PoolOptions{
			Dialer: func() (*beanstalk.Client, error) {
				return beanstalk.NewClient(mock.NewConn(nil, nil)), nil
			},
			Logger:      beanstalk.NopLogger,
			Capacity:    1,
			MaxAge:      0,
			IdleTimeout: 0,
		})

		require.Equal(t, 0, pool.Len())

		// opens pool
		require.NoError(t, pool.Open(context.Background()))

		require.Equal(t, 1, pool.Len())

		// gets client
		client, err := pool.Get()

		require.Nil(t, err)
		require.NotNil(t, client)

		require.NoError(t, client.Close())

		// puts client in pool
		require.NoError(t, pool.Put(client))

		require.Equal(t, 0, pool.Len())

		// closes pool
		require.NoError(t, pool.Close(context.Background()))

		require.Equal(t, 0, pool.Len())
	})

	t.Run("max capacity", func(t *testing.T) {
		pool := beanstalk.NewDefaultPool(&beanstalk.PoolOptions{
			Dialer: func() (*beanstalk.Client, error) {
				return beanstalk.NewClient(mock.NewConn(nil, nil)), nil
			},
			Logger:      beanstalk.NopLogger,
			Capacity:    1,
			MaxAge:      0,
			IdleTimeout: 0,
		})

		require.Equal(t, 0, pool.Len())

		// opens pool
		require.NoError(t, pool.Open(context.Background()))

		require.Equal(t, 1, pool.Len())

		// gets client
		client, err := pool.Get()

		require.Nil(t, err)
		require.NotNil(t, client)

		// waiting refill
		time.Sleep(1 * time.Second)

		// puts client in pool
		require.NoError(t, pool.Put(client))

		require.Equal(t, 1, pool.Len())

		// closes pool
		require.NoError(t, pool.Close(context.Background()))

		require.Equal(t, 0, pool.Len())
	})
}
