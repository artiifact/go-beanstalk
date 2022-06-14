package beanstalk

import (
	"gopkg.in/yaml.v2"
	"io"
	"net"
	"strings"
	"sync/atomic"
	"time"
)

type Client struct {
	conn      *Conn
	createdAt time.Time
	usedAt    int64
	closedAt  int64
}

func NewClient(conn io.ReadWriteCloser) *Client {
	return &Client{
		conn:      NewConn(conn),
		createdAt: time.Now(),
		usedAt:    0,
		closedAt:  0,
	}
}

func Dial(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return NewClient(conn), nil
}

func (c *Client) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Client) UsedAt() time.Time {
	return time.Unix(atomic.LoadInt64(&c.usedAt), 0)
}

func (c *Client) ClosedAt() time.Time {
	return time.Unix(atomic.LoadInt64(&c.closedAt), 0)
}

func (c *Client) Close() error {
	atomic.StoreInt64(&c.closedAt, time.Now().Unix())

	return c.conn.Close()
}

func (c *Client) Check() error {
	return c.conn.Check()
}

func (c *Client) Put(priority uint32, delay, ttr time.Duration, data []byte) (int, error) {
	r, err := c.ExecuteCommand(PutCommand{priority, delay, ttr, data})
	if err != nil {
		return 0, err
	}

	return r.(PutCommandResponse).ID, nil
}

func (c *Client) Use(tube string) (string, error) {
	r, err := c.ExecuteCommand(UseCommand{tube})
	if err != nil {
		return "", err
	}

	return r.(UseCommandResponse).Tube, nil
}

func (c *Client) Reserve() (Job, error) {
	r, err := c.ExecuteCommand(ReserveCommand{})
	if err != nil {
		return Job{}, err
	}

	return Job{r.(ReserveCommandResponse).ID, r.(ReserveCommandResponse).Data}, nil
}

func (c *Client) ReserveWithTimeout(timeout time.Duration) (Job, error) {
	r, err := c.ExecuteCommand(ReserveWithTimeoutCommand{timeout})
	if err != nil {
		return Job{}, err
	}

	return Job{r.(ReserveWithTimeoutCommandResponse).ID, r.(ReserveWithTimeoutCommandResponse).Data}, nil
}

func (c *Client) ReserveJob(id int) (Job, error) {
	r, err := c.ExecuteCommand(ReserveJobCommand{id})
	if err != nil {
		return Job{}, err
	}

	return Job{r.(ReserveJobCommandResponse).ID, r.(ReserveJobCommandResponse).Data}, nil
}

func (c *Client) Delete(id int) error {
	_, err := c.ExecuteCommand(DeleteCommand{id})

	return err
}

func (c *Client) Release(id int, priority uint32, delay time.Duration) error {
	_, err := c.ExecuteCommand(ReleaseCommand{id, priority, delay})

	return err
}

func (c *Client) Bury(id int, priority uint32) error {
	_, err := c.ExecuteCommand(BuryCommand{id, priority})

	return err
}

func (c *Client) Touch(id int) error {
	_, err := c.ExecuteCommand(TouchCommand{id})

	return err
}

func (c *Client) Watch(tube string) (int, error) {
	r, err := c.ExecuteCommand(WatchCommand{tube})
	if err != nil {
		return 0, err
	}

	return r.(WatchCommandResponse).Count, nil
}

func (c *Client) Ignore(tube string) (int, error) {
	r, err := c.ExecuteCommand(IgnoreCommand{tube})
	if err != nil {
		return 0, err
	}

	return r.(IgnoreCommandResponse).Count, nil
}

func (c *Client) Peek(id int) (Job, error) {
	r, err := c.ExecuteCommand(PeekCommand{id})
	if err != nil {
		return Job{}, err
	}

	return Job{r.(PeekCommandResponse).ID, r.(PeekCommandResponse).Data}, nil
}

func (c *Client) PeekReady() (Job, error) {
	r, err := c.ExecuteCommand(PeekReadyCommand{})
	if err != nil {
		return Job{}, err
	}

	return Job{r.(PeekReadyCommandResponse).ID, r.(PeekReadyCommandResponse).Data}, nil
}

func (c *Client) PeekDelayed() (Job, error) {
	r, err := c.ExecuteCommand(PeekDelayedCommand{})
	if err != nil {
		return Job{}, err
	}

	return Job{r.(PeekDelayedCommandResponse).ID, r.(PeekDelayedCommandResponse).Data}, nil
}

func (c *Client) PeekBuried() (Job, error) {
	r, err := c.ExecuteCommand(PeekBuriedCommand{})
	if err != nil {
		return Job{}, err
	}

	return Job{r.(PeekBuriedCommandResponse).ID, r.(PeekBuriedCommandResponse).Data}, nil
}

func (c *Client) Kick(bound int) (int, error) {
	r, err := c.ExecuteCommand(KickCommand{bound})
	if err != nil {
		return 0, err
	}

	return r.(KickCommandResponse).Count, nil
}

func (c *Client) KickJob(id int) error {
	_, err := c.ExecuteCommand(KickJobCommand{id})

	return err
}

func (c *Client) StatsJob(id int) (StatsJob, error) {
	r, err := c.ExecuteCommand(StatsJobCommand{id})
	if err != nil {
		return StatsJob{}, err
	}

	var stats StatsJob
	if err = yaml.Unmarshal(r.(StatsJobCommandResponse).Data, &stats); err != nil {
		return StatsJob{}, err
	}

	return stats, err
}

func (c *Client) StatsTube(tube string) (StatsTube, error) {
	r, err := c.ExecuteCommand(StatsTubeCommand{tube})
	if err != nil {
		return StatsTube{}, err
	}

	var stats StatsTube
	if err = yaml.Unmarshal(r.(StatsTubeCommandResponse).Data, &stats); err != nil {
		return StatsTube{}, err
	}

	return stats, err
}

func (c *Client) Stats() (Stats, error) {
	r, err := c.ExecuteCommand(StatsCommand{})
	if err != nil {
		return Stats{}, err
	}

	var stats Stats
	if err = yaml.Unmarshal(r.(StatsCommandResponse).Data, &stats); err != nil {
		return Stats{}, err
	}

	return stats, err
}

func (c *Client) ListTubes() ([]string, error) {
	r, err := c.ExecuteCommand(ListTubesCommand{})
	if err != nil {
		return nil, err
	}

	var tubes []string
	if err = yaml.Unmarshal(r.(ListTubesCommandResponse).Data, &tubes); err != nil {
		return nil, err
	}

	return tubes, nil
}

func (c *Client) ListTubeUsed() (string, error) {
	r, err := c.ExecuteCommand(ListTubeUsedCommand{})
	if err != nil {
		return "", err
	}

	return r.(ListTubeUsedCommandResponse).Tube, nil
}

func (c *Client) ListTubesWatched() ([]string, error) {
	r, err := c.ExecuteCommand(ListTubesWatchedCommand{})
	if err != nil {
		return nil, err
	}

	var tubes []string
	if err = yaml.Unmarshal(r.(ListTubesWatchedCommandResponse).Data, &tubes); err != nil {
		return nil, err
	}

	return tubes, nil
}

func (c *Client) PauseTube(tube string, delay time.Duration) error {
	_, err := c.ExecuteCommand(PauseTubeCommand{tube, delay})

	return err
}

func (c *Client) ExecuteCommand(command Command) (CommandResponse, error) {
	atomic.StoreInt64(&c.usedAt, time.Now().Unix())

	id, err := c.conn.WriteRequest(command.CommandLine(), command.Body())
	if err != nil {
		return nil, err
	}

	responseLine, body, err := c.conn.ReadResponse(id, command.HasResponseBody())
	if err != nil {
		return nil, err
	}

	switch {
	case strings.EqualFold(responseLine, "OUT_OF_MEMORY"):
		return nil, ErrOutOfMemory

	case strings.EqualFold(responseLine, "INTERNAL_ERROR"):
		return nil, ErrInternalError

	case strings.EqualFold(responseLine, "BAD_FORMAT"):
		return nil, ErrBadFormat

	case strings.EqualFold(responseLine, "UNKNOWN_COMMAND"):
		return nil, ErrUnknownCommand
	}

	if builder, ok := command.(CommandResponseBuilder); ok {
		return builder.BuildResponse(responseLine, body)
	}

	return nil, ErrMalformedCommand
}
