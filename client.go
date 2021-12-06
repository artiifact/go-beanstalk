package beanstalk

import (
	"gopkg.in/yaml.v2"
	"io"
	"net"
	"net/textproto"
	"strings"
	"time"
)

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
	ExecuteCommand(command Command) (CommandResponse, error)
}

type client struct {
	conn *Conn
}

func NewClient(conn io.ReadWriteCloser) Client {
	return &client{NewConn(textproto.NewConn(conn))}
}

func Dial(address string) (Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return NewClient(conn), nil
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) Put(priority uint32, delay, ttr time.Duration, data []byte) (int, error) {
	r, err := c.ExecuteCommand(PutCommand{priority, delay, ttr, data})
	if err != nil {
		return 0, err
	}

	return r.(PutCommandResponse).ID, nil
}

func (c *client) Use(tube string) (string, error) {
	r, err := c.ExecuteCommand(UseCommand{tube})
	if err != nil {
		return "", err
	}

	return r.(UseCommandResponse).Tube, nil
}

func (c *client) Reserve() (Job, error) {
	r, err := c.ExecuteCommand(ReserveCommand{})
	if err != nil {
		return Job{}, err
	}

	return Job{r.(ReserveCommandResponse).ID, r.(ReserveCommandResponse).Data}, nil
}

func (c *client) ReserveWithTimeout(timeout time.Duration) (Job, error) {
	r, err := c.ExecuteCommand(ReserveWithTimeoutCommand{timeout})
	if err != nil {
		return Job{}, err
	}

	return Job{r.(ReserveWithTimeoutCommandResponse).ID, r.(ReserveWithTimeoutCommandResponse).Data}, nil
}

func (c *client) ReserveJob(id int) (Job, error) {
	r, err := c.ExecuteCommand(ReserveJobCommand{id})
	if err != nil {
		return Job{}, err
	}

	return Job{r.(ReserveJobCommandResponse).ID, r.(ReserveJobCommandResponse).Data}, nil
}

func (c *client) Delete(id int) error {
	_, err := c.ExecuteCommand(DeleteCommand{id})

	return err
}

func (c *client) Release(id int, priority uint32, delay time.Duration) error {
	_, err := c.ExecuteCommand(ReleaseCommand{id, priority, delay})

	return err
}

func (c *client) Bury(id int, priority uint32) error {
	_, err := c.ExecuteCommand(BuryCommand{id, priority})

	return err
}

func (c *client) Touch(id int) error {
	_, err := c.ExecuteCommand(TouchCommand{id})

	return err
}

func (c *client) Watch(tube string) (int, error) {
	r, err := c.ExecuteCommand(WatchCommand{tube})
	if err != nil {
		return 0, err
	}

	return r.(WatchCommandResponse).Count, nil
}

func (c *client) Ignore(tube string) (int, error) {
	r, err := c.ExecuteCommand(IgnoreCommand{tube})
	if err != nil {
		return 0, err
	}

	return r.(IgnoreCommandResponse).Count, nil
}

func (c *client) Peek(id int) (Job, error) {
	r, err := c.ExecuteCommand(PeekCommand{id})
	if err != nil {
		return Job{}, err
	}

	return Job{r.(PeekCommandResponse).ID, r.(PeekCommandResponse).Data}, nil
}

func (c *client) PeekReady() (Job, error) {
	r, err := c.ExecuteCommand(PeekReadyCommand{})
	if err != nil {
		return Job{}, err
	}

	return Job{r.(PeekReadyCommandResponse).ID, r.(PeekReadyCommandResponse).Data}, nil
}

func (c *client) PeekDelayed() (Job, error) {
	r, err := c.ExecuteCommand(PeekDelayedCommand{})
	if err != nil {
		return Job{}, err
	}

	return Job{r.(PeekDelayedCommandResponse).ID, r.(PeekDelayedCommandResponse).Data}, nil
}

func (c *client) PeekBuried() (Job, error) {
	r, err := c.ExecuteCommand(PeekBuriedCommand{})
	if err != nil {
		return Job{}, err
	}

	return Job{r.(PeekBuriedCommandResponse).ID, r.(PeekBuriedCommandResponse).Data}, nil
}

func (c *client) Kick(bound int) (int, error) {
	r, err := c.ExecuteCommand(KickCommand{bound})
	if err != nil {
		return 0, err
	}

	return r.(KickCommandResponse).Count, nil
}

func (c *client) KickJob(id int) error {
	_, err := c.ExecuteCommand(KickJobCommand{id})

	return err
}

func (c *client) StatsJob(id int) (StatsJob, error) {
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

func (c *client) StatsTube(tube string) (StatsTube, error) {
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

func (c *client) Stats() (Stats, error) {
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

func (c *client) ListTubes() ([]string, error) {
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

func (c *client) ListTubeUsed() (string, error) {
	r, err := c.ExecuteCommand(ListTubeUsedCommand{})
	if err != nil {
		return "", err
	}

	return r.(ListTubeUsedCommandResponse).Tube, nil
}

func (c *client) ListTubesWatched() ([]string, error) {
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

func (c *client) PauseTube(tube string, delay time.Duration) error {
	_, err := c.ExecuteCommand(PauseTubeCommand{tube, delay})

	return err
}

func (c *client) ExecuteCommand(command Command) (CommandResponse, error) {
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

	return nil, nil
}
