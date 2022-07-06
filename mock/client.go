package mock

import (
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/stretchr/testify/mock"
	"time"
)

type Client struct {
	mock.Mock
}

func (c *Client) Close() error {
	args := c.Called()

	return args.Error(0)
}

func (c *Client) Put(priority uint32, delay, ttr time.Duration, data []byte) (int, error) {
	args := c.Called(priority, delay, ttr, data)

	return args.Int(0), args.Error(1)
}

func (c *Client) Use(tube string) (string, error) {
	args := c.Called(tube)

	return args.String(0), args.Error(1)
}

func (c *Client) Reserve() (*beanstalk.Job, error) {
	args := c.Called()

	return args.Get(0).(*beanstalk.Job), args.Error(1)
}

func (c *Client) ReserveWithTimeout(timeout time.Duration) (*beanstalk.Job, error) {
	args := c.Called(timeout)

	return args.Get(0).(*beanstalk.Job), args.Error(1)
}

func (c *Client) ReserveJob(id int) (*beanstalk.Job, error) {
	args := c.Called(id)

	return args.Get(0).(*beanstalk.Job), args.Error(1)
}

func (c *Client) Delete(id int) error {
	args := c.Called(id)

	return args.Error(0)
}

func (c *Client) Release(id int, priority uint32, delay time.Duration) error {
	args := c.Called(id, priority, delay)

	return args.Error(0)
}

func (c *Client) Bury(id int, priority uint32) error {
	args := c.Called(id, priority)

	return args.Error(0)
}

func (c *Client) Touch(id int) error {
	args := c.Called(id)

	return args.Error(0)
}

func (c *Client) Watch(tube string) (int, error) {
	args := c.Called(tube)

	return args.Int(0), args.Error(1)
}

func (c *Client) Ignore(tube string) (int, error) {
	args := c.Called(tube)

	return args.Int(0), args.Error(1)
}

func (c *Client) Peek(id int) (*beanstalk.Job, error) {
	args := c.Called(id)

	return args.Get(0).(*beanstalk.Job), args.Error(1)
}

func (c *Client) PeekReady() (*beanstalk.Job, error) {
	args := c.Called()

	return args.Get(0).(*beanstalk.Job), args.Error(1)
}

func (c *Client) PeekDelayed() (*beanstalk.Job, error) {
	args := c.Called()

	return args.Get(0).(*beanstalk.Job), args.Error(1)
}

func (c *Client) PeekBuried() (*beanstalk.Job, error) {
	args := c.Called()

	return args.Get(0).(*beanstalk.Job), args.Error(1)
}

func (c *Client) Kick(bound int) (int, error) {
	args := c.Called(bound)

	return args.Int(0), args.Error(1)
}

func (c *Client) KickJob(id int) error {
	args := c.Called(id)

	return args.Error(0)
}

func (c *Client) StatsJob(id int) (*beanstalk.StatsJob, error) {
	args := c.Called(id)

	return args.Get(0).(*beanstalk.StatsJob), args.Error(1)
}

func (c *Client) StatsTube(tube string) (*beanstalk.StatsTube, error) {
	args := c.Called(tube)

	return args.Get(0).(*beanstalk.StatsTube), args.Error(1)
}

func (c *Client) Stats() (*beanstalk.Stats, error) {
	args := c.Called()

	return args.Get(0).(*beanstalk.Stats), args.Error(1)
}

func (c *Client) ListTubes() ([]string, error) {
	args := c.Called()

	return args.Get(0).([]string), args.Error(1)
}

func (c *Client) ListTubeUsed() (string, error) {
	args := c.Called()

	return args.String(0), args.Error(1)
}

func (c *Client) ListTubesWatched() ([]string, error) {
	args := c.Called()

	return args.Get(0).([]string), args.Error(1)
}

func (c *Client) PauseTube(tube string, delay time.Duration) error {
	args := c.Called(tube, delay)

	return args.Error(0)
}
