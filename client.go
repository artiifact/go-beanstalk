package beanstalk

import (
	"github.com/IvanLutokhin/go-beanstalk/checker"
	"gopkg.in/yaml.v2"
	"io"
	"net"
	"net/textproto"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var crnl = []byte{'\r', '\n'}

type Client interface {
	Close() error
	Put(priority uint32, delay, ttr time.Duration, data []byte) (int, error)
	Use(tube string) (string, error)
	Reserve() (Job, error)
	ReserveWithTimeout(timeout time.Duration) (Job, error)
	ReserveJob(id int) (Job, error)
	Delete(id int) error
	Release(id int, priority uint32, delay time.Duration) error
	Bury(id int, priority uint32) error
	Touch(id int) error
	Watch(tube string) (int, error)
	Ignore(tube string) (int, error)
	Peek(id int) (Job, error)
	PeekReady() (Job, error)
	PeekDelayed() (Job, error)
	PeekBuried() (Job, error)
	Kick(bound int) (int, error)
	KickJob(id int) error
	StatsJob(id int) (StatsJob, error)
	StatsTube(tube string) (StatsTube, error)
	Stats() (Stats, error)
	ListTubes() ([]string, error)
	ListTubeUsed() (string, error)
	ListTubesWatched() ([]string, error)
	PauseTube(tube string, delay time.Duration) error
}

func Dial(address string) (Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return NewDefaultClient(conn), nil
}

type DefaultClient struct {
	conn      *textproto.Conn
	checker   *checker.Checker
	createdAt time.Time
	usedAt    int64
	closedAt  int64
}

func NewDefaultClient(conn io.ReadWriteCloser) *DefaultClient {
	return &DefaultClient{
		conn:      textproto.NewConn(conn),
		checker:   checker.New(conn),
		createdAt: time.Now(),
		usedAt:    0,
		closedAt:  0,
	}
}

func (c *DefaultClient) CreatedAt() time.Time {
	return c.createdAt
}

func (c *DefaultClient) UsedAt() time.Time {
	return time.Unix(atomic.LoadInt64(&c.usedAt), 0)
}

func (c *DefaultClient) ClosedAt() time.Time {
	return time.Unix(atomic.LoadInt64(&c.closedAt), 0)
}

func (c *DefaultClient) Check() error {
	return c.checker.Check()
}

func (c *DefaultClient) Close() error {
	atomic.StoreInt64(&c.closedAt, time.Now().Unix())

	return c.conn.Close()
}

func (c *DefaultClient) Put(priority uint32, delay, ttr time.Duration, data []byte) (int, error) {
	r, err := c.ExecuteCommand(PutCommand{Priority: priority, Delay: delay, TTR: ttr, Data: data})
	if err != nil {
		return 0, err
	}

	return r.(PutCommandResponse).ID, nil
}

func (c *DefaultClient) Use(tube string) (string, error) {
	r, err := c.ExecuteCommand(UseCommand{Tube: tube})
	if err != nil {
		return "", err
	}

	return r.(UseCommandResponse).Tube, nil
}

func (c *DefaultClient) Reserve() (Job, error) {
	r, err := c.ExecuteCommand(ReserveCommand{})
	if err != nil {
		return Job{}, err
	}

	return Job{ID: r.(ReserveCommandResponse).ID, Data: r.(ReserveCommandResponse).Data}, nil
}

func (c *DefaultClient) ReserveWithTimeout(timeout time.Duration) (Job, error) {
	r, err := c.ExecuteCommand(ReserveWithTimeoutCommand{Timeout: timeout})
	if err != nil {
		return Job{}, err
	}

	return Job{ID: r.(ReserveWithTimeoutCommandResponse).ID, Data: r.(ReserveWithTimeoutCommandResponse).Data}, nil
}

func (c *DefaultClient) ReserveJob(id int) (Job, error) {
	r, err := c.ExecuteCommand(ReserveJobCommand{ID: id})
	if err != nil {
		return Job{}, err
	}

	return Job{ID: r.(ReserveJobCommandResponse).ID, Data: r.(ReserveJobCommandResponse).Data}, nil
}

func (c *DefaultClient) Delete(id int) error {
	_, err := c.ExecuteCommand(DeleteCommand{ID: id})

	return err
}

func (c *DefaultClient) Release(id int, priority uint32, delay time.Duration) error {
	_, err := c.ExecuteCommand(ReleaseCommand{ID: id, Priority: priority, Delay: delay})

	return err
}

func (c *DefaultClient) Bury(id int, priority uint32) error {
	_, err := c.ExecuteCommand(BuryCommand{ID: id, Priority: priority})

	return err
}

func (c *DefaultClient) Touch(id int) error {
	_, err := c.ExecuteCommand(TouchCommand{ID: id})

	return err
}

func (c *DefaultClient) Watch(tube string) (int, error) {
	r, err := c.ExecuteCommand(WatchCommand{Tube: tube})
	if err != nil {
		return 0, err
	}

	return r.(WatchCommandResponse).Count, nil
}

func (c *DefaultClient) Ignore(tube string) (int, error) {
	r, err := c.ExecuteCommand(IgnoreCommand{Tube: tube})
	if err != nil {
		return 0, err
	}

	return r.(IgnoreCommandResponse).Count, nil
}

func (c *DefaultClient) Peek(id int) (Job, error) {
	r, err := c.ExecuteCommand(PeekCommand{ID: id})
	if err != nil {
		return Job{}, err
	}

	return Job{ID: r.(PeekCommandResponse).ID, Data: r.(PeekCommandResponse).Data}, nil
}

func (c *DefaultClient) PeekReady() (Job, error) {
	r, err := c.ExecuteCommand(PeekReadyCommand{})
	if err != nil {
		return Job{}, err
	}

	return Job{ID: r.(PeekReadyCommandResponse).ID, Data: r.(PeekReadyCommandResponse).Data}, nil
}

func (c *DefaultClient) PeekDelayed() (Job, error) {
	r, err := c.ExecuteCommand(PeekDelayedCommand{})
	if err != nil {
		return Job{}, err
	}

	return Job{ID: r.(PeekDelayedCommandResponse).ID, Data: r.(PeekDelayedCommandResponse).Data}, nil
}

func (c *DefaultClient) PeekBuried() (Job, error) {
	r, err := c.ExecuteCommand(PeekBuriedCommand{})
	if err != nil {
		return Job{}, err
	}

	return Job{ID: r.(PeekBuriedCommandResponse).ID, Data: r.(PeekBuriedCommandResponse).Data}, nil
}

func (c *DefaultClient) Kick(bound int) (int, error) {
	r, err := c.ExecuteCommand(KickCommand{Bound: bound})
	if err != nil {
		return 0, err
	}

	return r.(KickCommandResponse).Count, nil
}

func (c *DefaultClient) KickJob(id int) error {
	_, err := c.ExecuteCommand(KickJobCommand{ID: id})

	return err
}

func (c *DefaultClient) StatsJob(id int) (StatsJob, error) {
	r, err := c.ExecuteCommand(StatsJobCommand{ID: id})
	if err != nil {
		return StatsJob{}, err
	}

	var stats StatsJob
	if err = yaml.Unmarshal(r.(StatsJobCommandResponse).Data, &stats); err != nil {
		return StatsJob{}, err
	}

	return stats, err
}

func (c *DefaultClient) StatsTube(tube string) (StatsTube, error) {
	r, err := c.ExecuteCommand(StatsTubeCommand{Tube: tube})
	if err != nil {
		return StatsTube{}, err
	}

	var stats StatsTube
	if err = yaml.Unmarshal(r.(StatsTubeCommandResponse).Data, &stats); err != nil {
		return StatsTube{}, err
	}

	return stats, err
}

func (c *DefaultClient) Stats() (Stats, error) {
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

func (c *DefaultClient) ListTubes() ([]string, error) {
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

func (c *DefaultClient) ListTubeUsed() (string, error) {
	r, err := c.ExecuteCommand(ListTubeUsedCommand{})
	if err != nil {
		return "", err
	}

	return r.(ListTubeUsedCommandResponse).Tube, nil
}

func (c *DefaultClient) ListTubesWatched() ([]string, error) {
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

func (c *DefaultClient) PauseTube(tube string, delay time.Duration) error {
	_, err := c.ExecuteCommand(PauseTubeCommand{Tube: tube, Delay: delay})

	return err
}

func (c *DefaultClient) ExecuteCommand(command Command) (CommandResponse, error) {
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

func (c *DefaultClient) writeRequest(line string, body []byte) (uint, error) {
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

func (c *DefaultClient) readResponse(id uint, hasBody bool) (string, []byte, error) {
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
