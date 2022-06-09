package beanstalk_test

import (
	"context"
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk/pkg/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPHandlerAdapter_ServeHTTP(t *testing.T) {
	pool := beanstalk.NewPool(&beanstalk.PoolOptions{
		Dialer: func() (*beanstalk.Client, error) {
			return mock.NewClient(nil, nil), nil
		},
		Logger:      beanstalk.NopLogger,
		Capacity:    1,
		MaxAge:      0,
		IdleTimeout: 0,
	})

	handler := beanstalk.HandlerFunc(func(client *beanstalk.Client, writer http.ResponseWriter, request *http.Request) {
		require.NotNil(t, client)

		writer.WriteHeader(http.StatusOK)
	})

	adapter := beanstalk.NewHTTPHandlerAdapter(pool, handler)

	t.Run("panic on get client", func(t *testing.T) {
		defer func() {
			require.NotNil(t, recover())
		}()

		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		adapter.ServeHTTP(recorder, request)
	})

	t.Run("success", func(t *testing.T) {
		if err := pool.Open(context.Background()); err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		adapter.ServeHTTP(recorder, request)
	})

	if err := pool.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
}
