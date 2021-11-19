package beanstalk

import (
	"fmt"
	"strings"
	"time"
)

type ReleaseCommand struct {
	ID       int
	Priority uint32
	Delay    time.Duration
}

type ReleaseCommandResponse struct{}

func (c ReleaseCommand) CommandLine() string {
	return fmt.Sprintf("release %d %d %0.f", c.ID, c.Priority, c.Delay.Seconds())
}

func (c ReleaseCommand) Body() []byte {
	return nil
}

func (c ReleaseCommand) HasResponseBody() bool {
	return false
}

func (c ReleaseCommand) BuildResponse(responseLine string, _ []byte) (CommandResponse, error) {
	switch {
	case strings.EqualFold(responseLine, "RELEASED"):
		return ReleaseCommandResponse{}, nil

	case strings.EqualFold(responseLine, "BURIED"):
		return nil, ErrBuried

	case strings.EqualFold(responseLine, "NOT_FOUND"):
		return nil, ErrNotFound

	default:
		return nil, ErrUnexpectedResponse
	}
}
