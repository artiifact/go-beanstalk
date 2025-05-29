package mock_test

import (
	"io"
	"testing"

	"github.com/artiifact/go-beanstalk/mock"
	"github.com/stretchr/testify/require"
)

func TestConn_Read(t *testing.T) {
	t.Run("EOF / nil output", func(t *testing.T) {
		conn := mock.NewConn(nil, nil)

		b := make([]byte, 8)

		n, err := conn.Read(b)

		require.Equal(t, io.EOF, err)
		require.Equal(t, 0, n)
	})

	t.Run("EOF / empty output", func(t *testing.T) {
		conn := mock.NewConn(nil, []string{""})

		b := make([]byte, 8)

		n, err := conn.Read(b)

		require.Equal(t, io.EOF, err)
		require.Equal(t, 0, n)
	})

	t.Run("success", func(t *testing.T) {
		output := []string{"test 1", "test 2", "test 3"}

		conn := mock.NewConn(nil, output)

		for _, item := range output {
			b := make([]byte, 8)

			n, err := conn.Read(b)

			require.Nil(t, err)
			require.Equal(t, len(item), n)
			require.Equal(t, item, string(b[:n]))
		}

		require.NoError(t, conn.Close())
	})
}

func TestConn_Write(t *testing.T) {
	t.Run("EOF / nil input", func(t *testing.T) {
		conn := mock.NewConn(nil, nil)

		n, err := conn.Write([]byte{})

		require.Equal(t, io.EOF, err)
		require.Equal(t, 0, n)
	})

	t.Run("EOF / empty input", func(t *testing.T) {
		conn := mock.NewConn([]string{""}, nil)

		n, err := conn.Write([]byte{})

		require.Equal(t, io.EOF, err)
		require.Equal(t, 0, n)
	})

	t.Run("unexpected bytes", func(t *testing.T) {
		conn := mock.NewConn([]string{"test"}, nil)

		n, err := conn.Write([]byte("data"))

		require.NotNil(t, err)
		require.Equal(t, `beanstalk: conn: expected "test", got "data"`, err.Error())
		require.Equal(t, 0, n)
	})

	t.Run("success", func(t *testing.T) {
		input := []string{"test 1", "test 2", "test 3"}

		conn := mock.NewConn(input, nil)

		for _, item := range input {
			n, err := conn.Write([]byte(item))

			require.Nil(t, err)
			require.Equal(t, len(item), n)
		}

		require.NoError(t, conn.Close())
	})
}

func TestConn_Close(t *testing.T) {
	t.Run("input not empty", func(t *testing.T) {
		conn := mock.NewConn([]string{"test"}, nil)

		require.Error(t, conn.Close())
	})

	t.Run("output not empty", func(t *testing.T) {
		conn := mock.NewConn(nil, []string{"test"})

		require.Error(t, conn.Close())
	})

	t.Run("success", func(t *testing.T) {
		conn := mock.NewConn(nil, nil)

		require.NoError(t, conn.Close())
	})
}
