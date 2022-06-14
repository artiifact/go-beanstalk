package beanstalk_test

import (
	"context"
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk/pkg/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		pool := beanstalk.NewPool(&beanstalk.PoolOptions{
			Dialer: func() (*beanstalk.Client, error) {
				return mock.NewClient(nil, nil), nil
			},
		})

		require.Equal(t, 0, pool.Len())

		// opens pool
		if err := pool.Open(context.Background()); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, 1, pool.Len())

		// gets client from pool
		client, err := pool.Get()

		require.Nil(t, err)
		require.NotNil(t, client)

		// puts client in pool
		err = pool.Put(client)

		require.Nil(t, err)

		// closes pool
		if err := pool.Close(context.Background()); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, 0, pool.Len())
	})

	t.Run("dialer not specified", func(t *testing.T) {
		pool := beanstalk.NewPool(&beanstalk.PoolOptions{
			Logger:      beanstalk.NopLogger,
			Capacity:    3,
			MaxAge:      0,
			IdleTimeout: 0,
		})

		require.Equal(t, 0, pool.Len())

		// opens pool
		if err := pool.Open(context.Background()); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, 0, pool.Len())

		// gets client from pool
		client, err := pool.Get()

		require.Equal(t, beanstalk.ErrDialerNotSpecified, err)
		require.Nil(t, client)

		// closes pool
		if err := pool.Close(context.Background()); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, 0, pool.Len())
	})
}

func TestPool_Open(t *testing.T) {
	t.Run("timeout", func(t *testing.T) {
		pool := beanstalk.NewPool(&beanstalk.PoolOptions{
			Dialer: func() (*beanstalk.Client, error) {
				return mock.NewClient(nil, nil), nil
			},
			Logger:      beanstalk.NopLogger,
			Capacity:    3,
			MaxAge:      0,
			IdleTimeout: 0,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 0)

		defer cancel()

		err := pool.Open(ctx)

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "deadline exceeded")
	})

	t.Run("already opened", func(t *testing.T) {
		pool := beanstalk.NewPool(&beanstalk.PoolOptions{
			Dialer: func() (*beanstalk.Client, error) {
				return mock.NewClient(nil, nil), nil
			},
			Logger:      beanstalk.NopLogger,
			Capacity:    3,
			MaxAge:      0,
			IdleTimeout: 0,
		})

		require.Equal(t, 0, pool.Len())

		// opens pool
		if err := pool.Open(context.Background()); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, 3, pool.Len())

		// retries to open pool
		err := pool.Open(context.Background())

		require.Equal(t, beanstalk.ErrAlreadyOpenedPool, err)
	})
}

func TestPool_Close(t *testing.T) {
	t.Run("timeout", func(t *testing.T) {
		pool := beanstalk.NewPool(&beanstalk.PoolOptions{
			Dialer: func() (*beanstalk.Client, error) {
				return mock.NewClient(nil, nil), nil
			},
			Logger:      beanstalk.NopLogger,
			Capacity:    5,
			MaxAge:      0,
			IdleTimeout: 0,
		})

		require.Equal(t, 0, pool.Len())

		// opens pool
		if err := pool.Open(context.Background()); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, 5, pool.Len())

		// closes pool
		ctx, cancel := context.WithTimeout(context.Background(), 0)

		defer cancel()

		err := pool.Close(ctx)

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "deadline exceeded")
	})

	t.Run("closed", func(t *testing.T) {
		pool := beanstalk.NewPool(&beanstalk.PoolOptions{
			Dialer: func() (*beanstalk.Client, error) {
				return mock.NewClient(nil, nil), nil
			},
			Logger:      beanstalk.NopLogger,
			Capacity:    5,
			MaxAge:      0,
			IdleTimeout: 0,
		})

		require.Equal(t, 0, pool.Len())

		// opens pool
		if err := pool.Open(context.Background()); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, 5, pool.Len())

		// closes pool
		if err := pool.Close(context.Background()); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, 0, pool.Len())

		// retries to close pool
		err := pool.Close(context.Background())

		require.Equal(t, beanstalk.ErrClosedPool, err)
	})
}

func TestPool_Get(t *testing.T) {
	t.Run("max age", func(t *testing.T) {
		pool := beanstalk.NewPool(&beanstalk.PoolOptions{
			Dialer: func() (*beanstalk.Client, error) {
				return mock.NewClient(nil, nil), nil
			},
			Logger:      beanstalk.NopLogger,
			Capacity:    1,
			MaxAge:      1 * time.Second,
			IdleTimeout: 0,
		})

		require.Equal(t, 0, pool.Len())

		// opens pool
		if err := pool.Open(context.Background()); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, 1, pool.Len())

		// waiting timeout
		time.Sleep(1 * time.Second)

		// gets client by factory
		client, err := pool.Get()

		require.Nil(t, err)
		require.NotNil(t, client)

		// closes pool
		if err := pool.Close(context.Background()); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, 0, pool.Len())
	})

	t.Run("idle timeout", func(t *testing.T) {
		pool := beanstalk.NewPool(&beanstalk.PoolOptions{
			Dialer: func() (*beanstalk.Client, error) {
				return mock.NewClient(nil, nil), nil
			},
			Logger:      beanstalk.NopLogger,
			Capacity:    1,
			MaxAge:      0,
			IdleTimeout: 1 * time.Second,
		})

		require.Equal(t, 0, pool.Len())

		// opens pool
		if err := pool.Open(context.Background()); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, 1, pool.Len())

		// waiting timeout
		time.Sleep(1 * time.Second)

		// gets client by factory
		client, err := pool.Get()

		require.Nil(t, err)
		require.NotNil(t, client)

		// closes pool
		if err := pool.Close(context.Background()); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, 0, pool.Len())
	})

	t.Run("closed", func(t *testing.T) {
		pool := beanstalk.NewPool(&beanstalk.PoolOptions{
			Dialer: func() (*beanstalk.Client, error) {
				client := mock.NewClient(nil, nil)
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
		if err := pool.Open(context.Background()); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, 1, pool.Len())

		// gets client by factory
		client, err := pool.Get()

		require.Nil(t, err)
		require.NotNil(t, client)

		// closes pool
		if err := pool.Close(context.Background()); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, 0, pool.Len())
	})
}

func TestPool_Put(t *testing.T) {
	t.Run("closed pool", func(t *testing.T) {
		pool := beanstalk.NewPool(&beanstalk.PoolOptions{
			Dialer: func() (*beanstalk.Client, error) {
				return mock.NewClient(nil, nil), nil
			},
			Logger:      beanstalk.NopLogger,
			Capacity:    1,
			MaxAge:      0,
			IdleTimeout: 0,
		})

		require.Equal(t, 0, pool.Len())

		// opens pool
		if err := pool.Open(context.Background()); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, 1, pool.Len())

		// gets client
		client, err := pool.Get()

		require.Nil(t, err)
		require.NotNil(t, client)

		// closes pool
		if err := pool.Close(context.Background()); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, 0, pool.Len())

		// puts client in closed pool
		err = pool.Put(client)

		require.Equal(t, beanstalk.ErrClosedPool, err)
	})

	t.Run("closed client", func(t *testing.T) {
		pool := beanstalk.NewPool(&beanstalk.PoolOptions{
			Dialer: func() (*beanstalk.Client, error) {
				return mock.NewClient(nil, nil), nil
			},
			Logger:      beanstalk.NopLogger,
			Capacity:    1,
			MaxAge:      0,
			IdleTimeout: 0,
		})

		require.Equal(t, 0, pool.Len())

		// opens pool
		if err := pool.Open(context.Background()); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, 1, pool.Len())

		// gets client
		client, err := pool.Get()

		require.Nil(t, err)
		require.NotNil(t, client)

		if err := client.Close(); err != nil {
			t.Fatal(err)
		}

		// puts client in pool
		err = pool.Put(client)

		require.Nil(t, err)
		require.Equal(t, 0, pool.Len())

		// closes pool
		if err := pool.Close(context.Background()); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, 0, pool.Len())
	})

	t.Run("max capacity", func(t *testing.T) {
		pool := beanstalk.NewPool(&beanstalk.PoolOptions{
			Dialer: func() (*beanstalk.Client, error) {
				return mock.NewClient(nil, nil), nil
			},
			Logger:      beanstalk.NopLogger,
			Capacity:    1,
			MaxAge:      0,
			IdleTimeout: 0,
		})

		require.Equal(t, 0, pool.Len())

		// opens pool
		if err := pool.Open(context.Background()); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, 1, pool.Len())

		// gets client
		client, err := pool.Get()

		require.Nil(t, err)
		require.NotNil(t, client)

		// waiting refill
		time.Sleep(1 * time.Second)

		// puts client in pool
		err = pool.Put(client)

		require.Nil(t, err)
		require.Equal(t, 1, pool.Len())

		// closes pool
		if err := pool.Close(context.Background()); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, 0, pool.Len())
	})
}
