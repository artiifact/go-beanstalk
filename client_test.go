package beanstalk_test

import (
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk/mock"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
	"time"
)

func TestNewDefaultClient(t *testing.T) {
	c := beanstalk.NewDefaultClient(mock.NewConn(nil, nil))

	require.NotNil(t, c.CreatedAt())
	require.Equal(t, int64(0), c.UsedAt().Unix())
	require.Equal(t, int64(0), c.ClosedAt().Unix())

	require.NoError(t, c.Close())
}

func TestDefaultClient_Close(t *testing.T) {
	c := beanstalk.NewDefaultClient(mock.NewConn(nil, nil))

	require.Equal(t, int64(0), c.ClosedAt().Unix())

	require.NoError(t, c.Close())

	require.NotEqual(t, int64(0), c.ClosedAt().Unix())
}

func TestDefaultClient_Put(t *testing.T) {
	t.Run("inserted / unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"put 1 5 600 4\r\ntest\r\n"}, []string{"INSERTED\r\n"}))

		_, err := c.Put(1, 5*time.Second, 10*time.Minute, []byte("test"))

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)

		require.NoError(t, c.Close())
	})

	t.Run("inserted / malformed response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"put 1 5 600 4\r\ntest\r\n"}, []string{"INSERTED test\r\n"}))

		_, err := c.Put(1, 5*time.Second, 10*time.Minute, []byte("test"))

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "invalid syntax")

		require.NoError(t, c.Close())
	})

	t.Run("inserted / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"put 1 5 600 4\r\ntest\r\n"}, []string{"INSERTED 1\r\n"}))

		id, err := c.Put(1, 5*time.Second, 10*time.Minute, []byte("test"))

		require.Nil(t, err)
		require.Equal(t, 1, id)

		require.NoError(t, c.Close())
	})

	t.Run("buried", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"put 100 0 1800 11\r\ntest buried\r\n"}, []string{"BURIED 1\r\n"}))

		_, err := c.Put(100, 0, 30*time.Minute, []byte("test buried"))

		require.Equal(t, beanstalk.ErrBuried, err)

		require.NoError(t, c.Close())
	})

	t.Run("expected CRLF", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"put 0 30 90 18\r\ntest expected CRLF\r\n"}, []string{"EXPECTED_CRLF\r\n"}))

		_, err := c.Put(0, 30*time.Second, 90*time.Second, []byte("test expected CRLF"))

		require.Equal(t, beanstalk.ErrExpectedCRLF, err)

		require.NoError(t, c.Close())
	})

	t.Run("job too big", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"put 1 1 1 16\r\ntest job too big\r\n"}, []string{"JOB_TOO_BIG\r\n"}))

		_, err := c.Put(1, 1*time.Second, 1*time.Second, []byte("test job too big"))

		require.Equal(t, beanstalk.ErrJobTooBig, err)

		require.NoError(t, c.Close())
	})

	t.Run("draining", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"put 0 0 0 13\r\ntest draining\r\n"}, []string{"DRAINING\r\n"}))

		_, err := c.Put(0, 0, 0, []byte("test draining"))

		require.Equal(t, beanstalk.ErrDraining, err)

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"put 0 0 0 24\r\ntest unexpected response\r\n"}, []string{"TEST\r\n"}))

		_, err := c.Put(0, 0, 0, []byte("test unexpected response"))

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_Use(t *testing.T) {
	t.Run("using / unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"use test\r\n"}, []string{"USING\r\n"}))

		_, err := c.Use("test")

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)

		require.NoError(t, c.Close())
	})

	t.Run("using / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"use test\r\n"}, []string{"USING test\r\n"}))

		tube, err := c.Use("test")

		require.Nil(t, err)
		require.Equal(t, "test", tube)

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"use test\r\n"}, []string{"TEST\r\n"}))

		_, err := c.Use("test")

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_Reserve(t *testing.T) {
	t.Run("reserved / unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"reserve\r\n"}, []string{"RESERVED\r\n"}))

		job, err := c.Reserve()

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("reserved / malformed response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"reserve\r\n"}, []string{"RESERVED test 4\r\ntest\r\n"}))

		job, err := c.Reserve()

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "invalid syntax")
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("reserved / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"reserve\r\n"}, []string{"RESERVED 1 4\r\ntest\r\n"}))

		job, err := c.Reserve()

		require.Nil(t, err)
		require.Equal(t, 1, job.ID)
		require.Equal(t, []byte("test"), job.Data)

		require.NoError(t, c.Close())
	})

	t.Run("deadline soon", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"reserve\r\n"}, []string{"DEADLINE_SOON\r\n"}))

		job, err := c.Reserve()

		require.Equal(t, beanstalk.ErrDeadlineSoon, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("timed out", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"reserve\r\n"}, []string{"TIMED_OUT\r\n"}))

		job, err := c.Reserve()

		require.Equal(t, beanstalk.ErrTimedOut, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"reserve\r\n"}, []string{"TEST\r\n"}))

		job, err := c.Reserve()

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_ReserveWithTimeout(t *testing.T) {
	t.Run("reserved / unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"reserve-with-timeout 5\r\n"}, []string{"RESERVED\r\n"}))

		job, err := c.ReserveWithTimeout(5 * time.Second)

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("reserved / malformed response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"reserve-with-timeout 5\r\n"}, []string{"RESERVED test 4\r\ntest\r\n"}))

		job, err := c.ReserveWithTimeout(5 * time.Second)

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "invalid syntax")
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("reserved / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"reserve-with-timeout 5\r\n"}, []string{"RESERVED 1 4\r\ntest\r\n"}))

		job, err := c.ReserveWithTimeout(5 * time.Second)

		require.Nil(t, err)
		require.Equal(t, 1, job.ID)
		require.Equal(t, []byte("test"), job.Data)

		require.NoError(t, c.Close())
	})

	t.Run("deadline soon", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"reserve-with-timeout 60\r\n"}, []string{"DEADLINE_SOON\r\n"}))

		job, err := c.ReserveWithTimeout(60 * time.Second)

		require.Equal(t, beanstalk.ErrDeadlineSoon, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("timed out", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"reserve-with-timeout 600\r\n"}, []string{"TIMED_OUT\r\n"}))

		job, err := c.ReserveWithTimeout(10 * time.Minute)

		require.Equal(t, beanstalk.ErrTimedOut, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"reserve-with-timeout 300\r\n"}, []string{"TEST\r\n"}))

		job, err := c.ReserveWithTimeout(300 * time.Second)

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_ReserveJob(t *testing.T) {
	t.Run("reserved / unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"reserve-job 1\r\n"}, []string{"RESERVED\r\n"}))

		job, err := c.ReserveJob(1)

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("reserved / malformed response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"reserve-job 1\r\n"}, []string{"RESERVED test 4\r\ntest\r\n"}))

		job, err := c.ReserveJob(1)

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "invalid syntax")
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("reserved / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"reserve-job 1\r\n"}, []string{"RESERVED 1 4\r\ntest\r\n"}))

		job, err := c.ReserveJob(1)

		require.Nil(t, err)
		require.Equal(t, 1, job.ID)
		require.Equal(t, []byte("test"), job.Data)

		require.NoError(t, c.Close())
	})

	t.Run("not found", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"reserve-job 1\r\n"}, []string{"NOT_FOUND\r\n"}))

		job, err := c.ReserveJob(1)

		require.Equal(t, beanstalk.ErrNotFound, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"reserve-job 1\r\n"}, []string{"TEST\r\n"}))

		job, err := c.ReserveJob(1)

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_Delete(t *testing.T) {
	t.Run("deleted / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"delete 1\r\n"}, []string{"DELETED\r\n"}))

		err := c.Delete(1)

		require.Nil(t, err)

		require.NoError(t, c.Close())
	})

	t.Run("not found", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"delete 1\r\n"}, []string{"NOT_FOUND\r\n"}))

		err := c.Delete(1)

		require.Equal(t, beanstalk.ErrNotFound, err)

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"delete 1\r\n"}, []string{"TEST\r\n"}))

		err := c.Delete(1)

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_Release(t *testing.T) {
	t.Run("released / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"release 1 0 10\r\n"}, []string{"RELEASED\r\n"}))

		err := c.Release(1, 0, 10*time.Second)

		require.Nil(t, err)

		require.NoError(t, c.Close())
	})

	t.Run("buried", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"release 1 999 600\r\n"}, []string{"BURIED\r\n"}))

		err := c.Release(1, 999, 10*time.Minute)

		require.Equal(t, beanstalk.ErrBuried, err)

		require.NoError(t, c.Close())
	})

	t.Run("not found", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"release 1 10 600\r\n"}, []string{"NOT_FOUND\r\n"}))

		err := c.Release(1, 10, 10*time.Minute)

		require.Equal(t, beanstalk.ErrNotFound, err)

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"release 1 0 5\r\n"}, []string{"TEST\r\n"}))

		err := c.Release(1, 0, 5*time.Second)

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_Bury(t *testing.T) {
	t.Run("buried / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"bury 1 10\r\n"}, []string{"BURIED\r\n"}))

		err := c.Bury(1, 10)

		require.Nil(t, err)

		require.NoError(t, c.Close())
	})

	t.Run("not found", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"bury 999 0\r\n"}, []string{"NOT_FOUND\r\n"}))

		err := c.Bury(999, 0)

		require.Equal(t, beanstalk.ErrNotFound, err)

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"bury 1 0\r\n"}, []string{"TEST\r\n"}))

		err := c.Bury(1, 0)

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_Touch(t *testing.T) {
	t.Run("touched / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"touch 1\r\n"}, []string{"TOUCHED\r\n"}))

		err := c.Touch(1)

		require.Nil(t, err)

		require.NoError(t, c.Close())
	})

	t.Run("not found", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"touch 1\r\n"}, []string{"NOT_FOUND\r\n"}))

		err := c.Touch(1)

		require.Equal(t, beanstalk.ErrNotFound, err)

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"touch 1\r\n"}, []string{"TEST\r\n"}))

		err := c.Touch(1)

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_Watch(t *testing.T) {
	t.Run("watching / unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"watch test\r\n"}, []string{"WATCHING\r\n"}))

		_, err := c.Watch("test")

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)

		require.NoError(t, c.Close())
	})

	t.Run("watching / malformed response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"watch test\r\n"}, []string{"WATCHING test\r\n"}))

		_, err := c.Watch("test")

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "invalid syntax")

		require.NoError(t, c.Close())
	})

	t.Run("watching / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"watch test\r\n"}, []string{"WATCHING 1\r\n"}))

		count, err := c.Watch("test")

		require.Nil(t, err)
		require.Equal(t, 1, count)

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"watch test\r\n"}, []string{"TEST\r\n"}))

		_, err := c.Watch("test")

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_Ignore(t *testing.T) {
	t.Run("watching / unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"ignore test\r\n"}, []string{"WATCHING\r\n"}))

		_, err := c.Ignore("test")

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)

		require.NoError(t, c.Close())
	})

	t.Run("watching / malformed response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"ignore test\r\n"}, []string{"WATCHING test\r\n"}))

		_, err := c.Ignore("test")

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "invalid syntax")

		require.NoError(t, c.Close())
	})

	t.Run("watching / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"ignore test\r\n"}, []string{"WATCHING 1\r\n"}))

		count, err := c.Ignore("test")

		require.Nil(t, err)
		require.Equal(t, 1, count)

		require.NoError(t, c.Close())
	})

	t.Run("not ignored", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"ignore test\r\n"}, []string{"NOT_IGNORED\r\n"}))

		_, err := c.Ignore("test")

		require.Equal(t, beanstalk.ErrNotIgnored, err)

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"ignore test\r\n"}, []string{"TEST\r\n"}))

		_, err := c.Ignore("test")

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_Peek(t *testing.T) {
	t.Run("found / unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"peek 1\r\n"}, []string{"FOUND\r\n"}))

		job, err := c.Peek(1)

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("found / malformed response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"peek 1\r\n"}, []string{"FOUND test 4\r\ntest\r\n"}))

		job, err := c.Peek(1)

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "invalid syntax")
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("found / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"peek 1\r\n"}, []string{"FOUND 1 4\r\ntest\r\n"}))

		job, err := c.Peek(1)

		require.Nil(t, err)
		require.Equal(t, 1, job.ID)
		require.Equal(t, []byte("test"), job.Data)

		require.NoError(t, c.Close())
	})

	t.Run("not found", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"peek 1\r\n"}, []string{"NOT_FOUND\r\n"}))

		job, err := c.Peek(1)

		require.Equal(t, beanstalk.ErrNotFound, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"peek 1\r\n"}, []string{"TEST\r\n"}))

		job, err := c.Peek(1)

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_PeekReady(t *testing.T) {
	t.Run("found / unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"peek-ready\r\n"}, []string{"FOUND\r\n"}))

		job, err := c.PeekReady()

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("found / malformed response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"peek-ready\r\n"}, []string{"FOUND test 4\r\ntest\r\n"}))

		job, err := c.PeekReady()

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "invalid syntax")
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("found / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"peek-ready\r\n"}, []string{"FOUND 1 4\r\ntest\r\n"}))

		job, err := c.PeekReady()

		require.Nil(t, err)
		require.Equal(t, 1, job.ID)
		require.Equal(t, []byte("test"), job.Data)

		require.NoError(t, c.Close())
	})

	t.Run("not found", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"peek-ready\r\n"}, []string{"NOT_FOUND\r\n"}))

		job, err := c.PeekReady()

		require.Equal(t, beanstalk.ErrNotFound, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"peek-ready\r\n"}, []string{"TEST\r\n"}))

		job, err := c.PeekReady()

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_PeekDelayed(t *testing.T) {
	t.Run("found / unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"peek-delayed\r\n"}, []string{"FOUND\r\n"}))

		job, err := c.PeekDelayed()

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("found / malformed response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"peek-delayed\r\n"}, []string{"FOUND test 4\r\ntest\r\n"}))

		job, err := c.PeekDelayed()

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "invalid syntax")
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("found / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"peek-delayed\r\n"}, []string{"FOUND 1 4\r\ntest\r\n"}))

		job, err := c.PeekDelayed()

		require.Nil(t, err)
		require.Equal(t, 1, job.ID)
		require.Equal(t, []byte("test"), job.Data)

		require.NoError(t, c.Close())
	})

	t.Run("not found", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"peek-delayed\r\n"}, []string{"NOT_FOUND\r\n"}))

		job, err := c.PeekDelayed()

		require.Equal(t, beanstalk.ErrNotFound, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"peek-delayed\r\n"}, []string{"TEST\r\n"}))

		job, err := c.PeekDelayed()

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_PeekBuried(t *testing.T) {
	t.Run("found / unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"peek-buried\r\n"}, []string{"FOUND\r\n"}))

		job, err := c.PeekBuried()

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("found / malformed response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"peek-buried\r\n"}, []string{"FOUND test 4\r\ntest\r\n"}))

		job, err := c.PeekBuried()

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "invalid syntax")
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("found / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"peek-buried\r\n"}, []string{"FOUND 1 4\r\ntest\r\n"}))

		job, err := c.PeekBuried()

		require.Nil(t, err)
		require.Equal(t, 1, job.ID)
		require.Equal(t, []byte("test"), job.Data)

		require.NoError(t, c.Close())
	})

	t.Run("not found", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"peek-buried\r\n"}, []string{"NOT_FOUND\r\n"}))

		job, err := c.PeekBuried()

		require.Equal(t, beanstalk.ErrNotFound, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"peek-buried\r\n"}, []string{"TEST\r\n"}))

		job, err := c.PeekBuried()

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)
		require.Nil(t, job)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_Kick(t *testing.T) {
	t.Run("kicked / unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"kick 3\r\n"}, []string{"KICKED\r\n"}))

		_, err := c.Kick(3)

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)

		require.NoError(t, c.Close())
	})

	t.Run("kicked / malformed response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"kick 5\r\n"}, []string{"KICKED test\r\n"}))

		_, err := c.Kick(5)

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "invalid syntax")

		require.NoError(t, c.Close())
	})

	t.Run("kicked / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"kick 1\r\n"}, []string{"KICKED 1\r\n"}))

		count, err := c.Kick(1)

		require.Nil(t, err)
		require.Equal(t, 1, count)

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"kick 10\r\n"}, []string{"TEST\r\n"}))

		_, err := c.Kick(10)

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_KickJob(t *testing.T) {
	t.Run("kicked / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"kick-job 1\r\n"}, []string{"KICKED\r\n"}))

		err := c.KickJob(1)

		require.Nil(t, err)

		require.NoError(t, c.Close())
	})

	t.Run("not found", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"kick-job 1\r\n"}, []string{"NOT_FOUND\r\n"}))

		err := c.KickJob(1)

		require.Equal(t, beanstalk.ErrNotFound, err)

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"kick-job 1\r\n"}, []string{"TEST\r\n"}))

		err := c.KickJob(1)

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_StatsJob(t *testing.T) {
	t.Run("ok / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn(
			[]string{"stats-job 1\r\n"},
			[]string{
				"OK 148\r\n" +
					"---\n" +
					"id: 1\n" +
					"tube: default\n" +
					"state: ready\n" +
					"pri: 999\n" +
					"age: 12\n" +
					"delay: 15\n" +
					"ttr: 1\n" +
					"time-left: 10\n" +
					"file: 1\n" +
					"reserves: 1\n" +
					"timeouts: 1\n" +
					"releases: 1\n" +
					"buries: 1\n" +
					"kicks: 1\n" +
					"\r\n",
			},
		))

		stats, err := c.StatsJob(1)

		require.Nil(t, err)
		require.Equal(t, 1, stats.ID)
		require.Equal(t, "default", stats.Tube)
		require.Equal(t, "ready", stats.State)
		require.Equal(t, 999, stats.Priority)
		require.Equal(t, 12, stats.Age)
		require.Equal(t, 15, stats.Delay)
		require.Equal(t, 1, stats.TTR)
		require.Equal(t, 10, stats.TimeLeft)
		require.Equal(t, 1, stats.File)
		require.Equal(t, 1, stats.Reserves)
		require.Equal(t, 1, stats.Timeouts)
		require.Equal(t, 1, stats.Releases)
		require.Equal(t, 1, stats.Buries)
		require.Equal(t, 1, stats.Kicks)

		require.NoError(t, c.Close())
	})

	t.Run("ok / failure", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn(
			[]string{"stats-job 1\r\n"},
			[]string{
				"OK 6\r\n" +
					"---\n" +
					"test\n" +
					"\r\n",
			},
		))

		stats, err := c.StatsJob(1)

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "cannot unmarshal")
		require.Nil(t, stats)

		require.NoError(t, c.Close())
	})

	t.Run("not found", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"stats-job 1\r\n"}, []string{"NOT_FOUND\r\n"}))

		stats, err := c.StatsJob(1)

		require.Equal(t, beanstalk.ErrNotFound, err)
		require.Nil(t, stats)

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"stats-job 1\r\n"}, []string{"TEST\r\n"}))

		stats, err := c.StatsJob(1)

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)
		require.Nil(t, stats)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_StatsTube(t *testing.T) {
	t.Run("ok / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn(
			[]string{"stats-tube default\r\n"},
			[]string{
				"OK 268\r\n" +
					"---\n" +
					"name: default\n" +
					"current-jobs-urgent: 1\n" +
					"current-jobs-ready: 1\n" +
					"current-jobs-reserved: 1\n" +
					"current-jobs-delayed: 1\n" +
					"current-jobs-buried: 1\n" +
					"total-jobs: 5\n" +
					"current-using: 3\n" +
					"current-watching: 3\n" +
					"current-waiting: 2\n" +
					"cmd-delete: 1\n" +
					"cmd-pause-tube: 1\n" +
					"pause: 100\n" +
					"pause-time-left: 10\n" +
					"\r\n",
			},
		))

		stats, err := c.StatsTube("default")

		require.Nil(t, err)
		require.Equal(t, "default", stats.Name)
		require.Equal(t, 1, stats.CurrentJobsUrgent)
		require.Equal(t, 1, stats.CurrentJobsReady)
		require.Equal(t, 1, stats.CurrentJobsReserved)
		require.Equal(t, 1, stats.CurrentJobsDelayed)
		require.Equal(t, 1, stats.CurrentJobsBuried)
		require.Equal(t, 5, stats.TotalJobs)
		require.Equal(t, 3, stats.CurrentUsing)
		require.Equal(t, 3, stats.CurrentWatching)
		require.Equal(t, 2, stats.CurrentWaiting)
		require.Equal(t, 1, stats.CmdDelete)
		require.Equal(t, 1, stats.CmdPauseTube)
		require.Equal(t, 100, stats.Pause)
		require.Equal(t, 10, stats.PauseTimeLeft)

		require.NoError(t, c.Close())
	})

	t.Run("ok / failure", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn(
			[]string{"stats-tube test\r\n"},
			[]string{
				"OK 6\r\n" +
					"---\n" +
					"test\n" +
					"\r\n",
			},
		))

		stats, err := c.StatsTube("test")

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "cannot unmarshal")
		require.Nil(t, stats)

		require.NoError(t, c.Close())
	})

	t.Run("not found", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"stats-tube test\r\n"}, []string{"NOT_FOUND\r\n"}))

		stats, err := c.StatsTube("test")

		require.Equal(t, beanstalk.ErrNotFound, err)
		require.Nil(t, stats)

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"stats-tube test\r\n"}, []string{"TEST\r\n"}))

		stats, err := c.StatsTube("test")

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)
		require.Nil(t, stats)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_Stats(t *testing.T) {
	t.Run("ok / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn(
			[]string{"stats\r\n"},
			[]string{
				"OK 874\r\n" +
					"---\n" +
					"current-jobs-urgent: 1\n" +
					"current-jobs-ready: 1\n" +
					"current-jobs-reserved: 1\n" +
					"current-jobs-delayed: 1\n" +
					"current-jobs-buried: 1\n" +
					"cmd-put: 1\n" +
					"cmd-peek: 1\n" +
					"cmd-peek-ready: 1\n" +
					"cmd-peek-delayed: 1\n" +
					"cmd-peek-buried: 1\n" +
					"cmd-reserve: 1\n" +
					"cmd-use: 1\n" +
					"cmd-watch: 1\n" +
					"cmd-ignore: 1\n" +
					"cmd-delete: 1\n" +
					"cmd-release: 1\n" +
					"cmd-bury: 1\n" +
					"cmd-kick: 1\n" +
					"cmd-touch: 1\n" +
					"cmd-stats: 1\n" +
					"cmd-stats-job: 1\n" +
					"cmd-stats-tube: 1\n" +
					"cmd-list-tubes: 1\n" +
					"cmd-list-tube-used: 1\n" +
					"cmd-list-tubes-watched: 1\n" +
					"cmd-pause-tube: 1\n" +
					"job-timeouts: 10\n" +
					"total-jobs: 25\n" +
					"max-job-size: 65535\n" +
					"current-tubes: 1\n" +
					"current-connections: 3\n" +
					"current-producers: 2\n" +
					"current-workers: 5\n" +
					"current-waiting: 1\n" +
					"total-connections: 3\n" +
					"pid: 1\n" +
					"version: 1.10\n" +
					"rusage-utime: 0.148125\n" +
					"rusage-stime: 0.014812\n" +
					"uptime: 1864\n" +
					"binlog-oldest-index: 1\n" +
					"binlog-current-index: 1\n" +
					"binlog-records-migrated: 1\n" +
					"binlog-records-written: 1\n" +
					"binlog-max-size: 10485760\n" +
					"id: f40521014b63360d\n" +
					"hostname: 671db3de0474\n" +
					"\r\n",
			},
		))

		stats, err := c.Stats()

		require.Nil(t, err)
		require.Equal(t, 1, stats.CurrentJobsUrgent)
		require.Equal(t, 1, stats.CurrentJobsReady)
		require.Equal(t, 1, stats.CurrentJobsReserved)
		require.Equal(t, 1, stats.CurrentJobsDelayed)
		require.Equal(t, 1, stats.CurrentJobsBuried)
		require.Equal(t, 1, stats.CmdPut)
		require.Equal(t, 1, stats.CmdPeek)
		require.Equal(t, 1, stats.CmdPeekReady)
		require.Equal(t, 1, stats.CmdPeekDelayed)
		require.Equal(t, 1, stats.CmdPeekBuried)
		require.Equal(t, 1, stats.CmdReserve)
		require.Equal(t, 1, stats.CmdUse)
		require.Equal(t, 1, stats.CmdWatch)
		require.Equal(t, 1, stats.CmdIgnore)
		require.Equal(t, 1, stats.CmdDelete)
		require.Equal(t, 1, stats.CmdRelease)
		require.Equal(t, 1, stats.CmdBury)
		require.Equal(t, 1, stats.CmdKick)
		require.Equal(t, 1, stats.CmdTouch)
		require.Equal(t, 1, stats.CmdStats)
		require.Equal(t, 1, stats.CmdStatsJob)
		require.Equal(t, 1, stats.CmdStatsTube)
		require.Equal(t, 1, stats.CmdListTubes)
		require.Equal(t, 1, stats.CmdListTubeUsed)
		require.Equal(t, 1, stats.CmdListTubesWatched)
		require.Equal(t, 1, stats.CmdPauseTube)
		require.Equal(t, 10, stats.JobTimeouts)
		require.Equal(t, 25, stats.TotalJobs)
		require.Equal(t, 65535, stats.MaxJobSize)
		require.Equal(t, 1, stats.CurrentTubes)
		require.Equal(t, 3, stats.CurrentConnections)
		require.Equal(t, 2, stats.CurrentProducers)
		require.Equal(t, 5, stats.CurrentWorkers)
		require.Equal(t, 1, stats.CurrentWaiting)
		require.Equal(t, 3, stats.TotalConnections)
		require.Equal(t, 1, stats.PID)
		require.Equal(t, "1.10", stats.Version)
		require.Equal(t, 0.148125, stats.RUsageUTime)
		require.Equal(t, 0.014812, stats.RUsageSTime)
		require.Equal(t, 1864, stats.Uptime)
		require.Equal(t, 1, stats.BinlogOldestIndex)
		require.Equal(t, 1, stats.BinlogCurrentIndex)
		require.Equal(t, 1, stats.BinlogRecordsMigrated)
		require.Equal(t, 1, stats.BinlogRecordsWritten)
		require.Equal(t, 10485760, stats.BinlogMaxSize)
		require.Equal(t, "f40521014b63360d", stats.ID)
		require.Equal(t, "671db3de0474", stats.Hostname)

		require.NoError(t, c.Close())
	})

	t.Run("ok / failure", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn(
			[]string{"stats\r\n"},
			[]string{
				"OK 6\r\n" +
					"---\n" +
					"test\n" +
					"\r\n",
			},
		))

		stats, err := c.Stats()

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "cannot unmarshal")
		require.Nil(t, stats)

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"stats\r\n"}, []string{"TEST\r\n"}))

		stats, err := c.Stats()

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)
		require.Nil(t, stats)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_ListTubes(t *testing.T) {
	t.Run("ok / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"list-tubes\r\n"}, []string{"OK 21\r\n---\n- default\n- test\n\r\n"}))

		tubes, err := c.ListTubes()

		require.Nil(t, err)
		require.ElementsMatch(t, []string{"default", "test"}, tubes)

		if err := c.Close(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ok / failure", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn(
			[]string{"list-tubes\r\n"},
			[]string{
				"OK 6\r\n" +
					"---\n" +
					"test\n" +
					"\r\n",
			},
		))

		_, err := c.ListTubes()

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "cannot unmarshal")

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"list-tubes\r\n"}, []string{"TEST\r\n"}))

		_, err := c.ListTubes()

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_ListTubeUsed(t *testing.T) {
	t.Run("using / unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"list-tube-used\r\n"}, []string{"USING\r\n"}))

		_, err := c.ListTubeUsed()

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)

		require.NoError(t, c.Close())
	})

	t.Run("using / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"list-tube-used\r\n"}, []string{"USING test\r\n"}))

		tube, err := c.ListTubeUsed()

		require.Nil(t, err)
		require.Equal(t, "test", tube)

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"list-tube-used\r\n"}, []string{"TEST\r\n"}))

		_, err := c.ListTubeUsed()

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_ListTubesWatched(t *testing.T) {
	t.Run("ok / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"list-tubes-watched\r\n"}, []string{"OK 14\r\n---\n- default\n\r\n"}))

		tubes, err := c.ListTubesWatched()

		require.Nil(t, err)
		require.ElementsMatch(t, []string{"default"}, tubes)

		require.NoError(t, c.Close())
	})

	t.Run("ok / failure", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn(
			[]string{"list-tubes-watched\r\n"},
			[]string{
				"OK 6\r\n" +
					"---\n" +
					"test\n" +
					"\r\n",
			},
		))

		_, err := c.ListTubesWatched()

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "cannot unmarshal")

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"list-tubes-watched\r\n"}, []string{"TEST\r\n"}))

		_, err := c.ListTubesWatched()

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_PauseTube(t *testing.T) {
	t.Run("paused / success", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"pause-tube test 60\r\n"}, []string{"PAUSED\r\n"}))

		err := c.PauseTube("test", 60*time.Second)

		require.Nil(t, err)

		require.NoError(t, c.Close())
	})

	t.Run("not found", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"pause-tube test 10\r\n"}, []string{"NOT_FOUND\r\n"}))

		err := c.PauseTube("test", 10*time.Second)

		require.Equal(t, beanstalk.ErrNotFound, err)

		require.NoError(t, c.Close())
	})

	t.Run("unexpected response", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"pause-tube test 0\r\n"}, []string{"TEST\r\n"}))

		err := c.PauseTube("test", 0)

		require.Equal(t, beanstalk.ErrUnexpectedResponse, err)

		require.NoError(t, c.Close())
	})
}

func TestDefaultClient_ExecuteCommand(t *testing.T) {
	t.Run("write failure", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn(nil, nil))

		_, err := c.ExecuteCommand(mockCommand{})

		require.Equal(t, io.EOF, err)

		require.NoError(t, c.Close())
	})

	t.Run("read failure", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"mock\r\n"}, nil))

		_, err := c.ExecuteCommand(mockCommand{})

		require.Equal(t, io.EOF, err)

		require.NoError(t, c.Close())
	})

	t.Run("out of memory", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"mock\r\n"}, []string{"OUT_OF_MEMORY\r\n"}))

		_, err := c.ExecuteCommand(mockCommand{})

		require.Equal(t, beanstalk.ErrOutOfMemory, err)

		require.NoError(t, c.Close())
	})

	t.Run("internal error", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"mock\r\n"}, []string{"INTERNAL_ERROR\r\n"}))

		_, err := c.ExecuteCommand(mockCommand{})

		require.Equal(t, beanstalk.ErrInternalError, err)

		require.NoError(t, c.Close())
	})

	t.Run("bad format", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"mock\r\n"}, []string{"BAD_FORMAT\r\n"}))

		_, err := c.ExecuteCommand(mockCommand{})

		require.Equal(t, beanstalk.ErrBadFormat, err)

		require.NoError(t, c.Close())
	})

	t.Run("unknown command", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"mock\r\n"}, []string{"UNKNOWN_COMMAND\r\n"}))

		_, err := c.ExecuteCommand(mockCommand{})

		require.Equal(t, beanstalk.ErrUnknownCommand, err)

		require.NoError(t, c.Close())
	})

	t.Run("malformed command", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"mock\r\n"}, []string{"MALFORMED\r\n"}))

		_, err := c.ExecuteCommand(mockCommand{})

		require.Equal(t, beanstalk.ErrMalformedCommand, err)

		require.NoError(t, c.Close())
	})

	t.Run("use", func(t *testing.T) {
		c := beanstalk.NewDefaultClient(mock.NewConn([]string{"mock\r\n"}, []string{"OK\r\n"}))

		require.Equal(t, int64(0), c.UsedAt().Unix())

		_, _ = c.ExecuteCommand(mockCommand{})

		require.NotEqual(t, int64(0), c.UsedAt().Unix())

		require.NoError(t, c.Close())
	})
}

// mock command

type mockCommand struct{}

func (c mockCommand) CommandLine() string {
	return "mock"
}

func (c mockCommand) Body() []byte {
	return nil
}

func (c mockCommand) HasResponseBody() bool {
	return false
}
