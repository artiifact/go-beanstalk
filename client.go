package beanstalk

import (
	"io"
	"net"
	"net/textproto"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/artiifact/go-beanstalk/checker"
	"gopkg.in/yaml.v2"
)

var crnl = []byte{'\r', '\n'}

type Client struct {
	conn      *textproto.Conn
	checker   *checker.Checker
	createdAt time.Time
	usedAt    int64
	closedAt  int64
}

func Dial(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return NewClient(conn), nil
}

func NewClient(conn io.ReadWriteCloser) *Client {
	return &Client{
		conn:      textproto.NewConn(conn),
		checker:   checker.New(conn),
		createdAt: time.Now(),
		usedAt:    0,
		closedAt:  0,
	}
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

func (c *Client) Check() error {
	return c.checker.Check()
}

func (c *Client) Close() error {
	atomic.StoreInt64(&c.closedAt, time.Now().Unix())

	return c.conn.Close()
}

func (c *Client) Put(priority uint32, delay, ttr time.Duration, data []byte) (int, error) {
	r, err := c.ExecuteCommand(PutCommand{Priority: priority, Delay: delay, TTR: ttr, Data: data})
	if err != nil {
		return 0, err
	}

	return r.(PutCommandResponse).ID, nil
}

func (c *Client) Use(tube string) (string, error) {
	r, err := c.ExecuteCommand(UseCommand{Tube: tube})
	if err != nil {
		return "", err
	}

	return r.(UseCommandResponse).Tube, nil
}

func (c *Client) Reserve() (*Job, error) {
	r, err := c.ExecuteCommand(ReserveCommand{})
	if err != nil {
		return nil, err
	}

	return &Job{ID: r.(ReserveCommandResponse).ID, Data: r.(ReserveCommandResponse).Data}, nil
}

func (c *Client) ReserveWithTimeout(timeout time.Duration) (*Job, error) {
	r, err := c.ExecuteCommand(ReserveWithTimeoutCommand{Timeout: timeout})
	if err != nil {
		return nil, err
	}

	return &Job{ID: r.(ReserveWithTimeoutCommandResponse).ID, Data: r.(ReserveWithTimeoutCommandResponse).Data}, nil
}

func (c *Client) ReserveJob(id int) (*Job, error) {
	r, err := c.ExecuteCommand(ReserveJobCommand{ID: id})
	if err != nil {
		return nil, err
	}

	return &Job{ID: r.(ReserveJobCommandResponse).ID, Data: r.(ReserveJobCommandResponse).Data}, nil
}

func (c *Client) Delete(id int) error {
	_, err := c.ExecuteCommand(DeleteCommand{ID: id})

	return err
}

func (c *Client) Release(id int, priority uint32, delay time.Duration) error {
	_, err := c.ExecuteCommand(ReleaseCommand{ID: id, Priority: priority, Delay: delay})

	return err
}

func (c *Client) Bury(id int, priority uint32) error {
	_, err := c.ExecuteCommand(BuryCommand{ID: id, Priority: priority})

	return err
}

func (c *Client) Touch(id int) error {
	_, err := c.ExecuteCommand(TouchCommand{ID: id})

	return err
}

func (c *Client) Watch(tube string) (int, error) {
	r, err := c.ExecuteCommand(WatchCommand{Tube: tube})
	if err != nil {
		return 0, err
	}

	return r.(WatchCommandResponse).Count, nil
}

func (c *Client) Ignore(tube string) (int, error) {
	r, err := c.ExecuteCommand(IgnoreCommand{Tube: tube})
	if err != nil {
		return 0, err
	}

	return r.(IgnoreCommandResponse).Count, nil
}

func (c *Client) Peek(id int) (*Job, error) {
	r, err := c.ExecuteCommand(PeekCommand{ID: id})
	if err != nil {
		return nil, err
	}

	return &Job{ID: r.(PeekCommandResponse).ID, Data: r.(PeekCommandResponse).Data}, nil
}

func (c *Client) PeekReady() (*Job, error) {
	r, err := c.ExecuteCommand(PeekReadyCommand{})
	if err != nil {
		return nil, err
	}

	return &Job{ID: r.(PeekReadyCommandResponse).ID, Data: r.(PeekReadyCommandResponse).Data}, nil
}

func (c *Client) PeekDelayed() (*Job, error) {
	r, err := c.ExecuteCommand(PeekDelayedCommand{})
	if err != nil {
		return nil, err
	}

	return &Job{ID: r.(PeekDelayedCommandResponse).ID, Data: r.(PeekDelayedCommandResponse).Data}, nil
}

func (c *Client) PeekBuried() (*Job, error) {
	r, err := c.ExecuteCommand(PeekBuriedCommand{})
	if err != nil {
		return nil, err
	}

	return &Job{ID: r.(PeekBuriedCommandResponse).ID, Data: r.(PeekBuriedCommandResponse).Data}, nil
}

func (c *Client) Kick(bound int) (int, error) {
	r, err := c.ExecuteCommand(KickCommand{Bound: bound})
	if err != nil {
		return 0, err
	}

	return r.(KickCommandResponse).Count, nil
}

func (c *Client) KickJob(id int) error {
	_, err := c.ExecuteCommand(KickJobCommand{ID: id})

	return err
}

func (c *Client) StatsJob(id int) (*StatsJob, error) {
	r, err := c.ExecuteCommand(StatsJobCommand{ID: id})
	if err != nil {
		return nil, err
	}

	var stats StatsJob
	if err = yaml.Unmarshal(r.(StatsJobCommandResponse).Data, &stats); err != nil {
		return nil, err
	}

	return &stats, err
}

func (c *Client) StatsTube(tube string) (*StatsTube, error) {
	r, err := c.ExecuteCommand(StatsTubeCommand{Tube: tube})
	if err != nil {
		return nil, err
	}

	var stats StatsTube
	if err = yaml.Unmarshal(r.(StatsTubeCommandResponse).Data, &stats); err != nil {
		return nil, err
	}

	return &stats, err
}

func (c *Client) Stats() (*Stats, error) {
	r, err := c.ExecuteCommand(StatsCommand{})
	if err != nil {
		return nil, err
	}

	var stats Stats
	if err = yaml.Unmarshal(r.(StatsCommandResponse).Data, &stats); err != nil {
		return nil, err
	}

	return &stats, err
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
	_, err := c.ExecuteCommand(PauseTubeCommand{Tube: tube, Delay: delay})

	return err
}

func (c *Client) ExecuteCommand(command Command) (CommandResponse, error) {
	atomic.StoreInt64(&c.usedAt, time.Now().Unix())

	id, err := c.writeRequest(command.CommandLine(), command.Body())
	if err != nil {
		return nil, err
	}

	responseLine, body, err := c.readResponse(id, command.HasResponseBody())
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

func (c *Client) writeRequest(line string, body []byte) (uint, error) {
	id := c.conn.Next()

	c.conn.StartRequest(id)
	defer c.conn.EndRequest(id)

	if _, err := c.conn.W.Write([]byte(line)); err != nil {
		return 0, err
	}

	if _, err := c.conn.W.Write(crnl); err != nil {
		return 0, err
	}

	if body != nil {
		if _, err := c.conn.W.Write(body); err != nil {
			return 0, err
		}

		if _, err := c.conn.W.Write(crnl); err != nil {
			return 0, err
		}
	}

	if err := c.conn.W.Flush(); err != nil {
		return 0, err
	}

	return id, nil
}

func (c *Client) readResponse(id uint, hasBody bool) (string, []byte, error) {
	c.conn.StartResponse(id)
	defer c.conn.EndResponse(id)

	line, err := c.conn.ReadLine()
	if err != nil {
		return line, nil, err
	}

	var body []byte
	if hasBody {
		i := strings.LastIndex(line, " ")
		if i == -1 {
			return line, body, nil
		}

		n, err := strconv.Atoi(line[i+1:])
		if err != nil {
			return line, body, err
		}

		body = make([]byte, n+2) // include CR NL
		if _, err = io.ReadFull(c.conn.R, body); err != nil {
			return line, body, err
		}

		body = body[:n] // exclude CR NL
	}

	return line, body, nil
}
