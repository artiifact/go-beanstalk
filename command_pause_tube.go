package beanstalk

import (
	"fmt"
	"strings"
	"time"
)

type PauseTubeCommand struct {
	Tube  string
	Delay time.Duration
}

type PauseTubeCommandResponse struct{}

func (c PauseTubeCommand) CommandLine() string {
	return fmt.Sprintf("pause-tube %s %0.f", c.Tube, c.Delay.Seconds())
}

func (c PauseTubeCommand) Body() []byte {
	return nil
}

func (c PauseTubeCommand) HasResponseBody() bool {
	return false
}

func (c PauseTubeCommand) BuildResponse(responseLine string, _ []byte) (CommandResponse, error) {
	switch {
	case strings.EqualFold(responseLine, "PAUSED"):
		return PauseTubeCommandResponse{}, nil

	case strings.EqualFold(responseLine, "NOT_FOUND"):
		return nil, ErrNotFound

	default:
		return nil, ErrUnexpectedResponse
	}
}
