package beanstalk

import (
	"fmt"
	"strings"
)

type StatsJobCommand struct {
	ID int
}

type StatsJobCommandResponse struct {
	Data []byte
}

func (c StatsJobCommand) CommandLine() string {
	return fmt.Sprintf("stats-job %d", c.ID)
}

func (c StatsJobCommand) Body() []byte {
	return nil
}

func (c StatsJobCommand) HasResponseBody() bool {
	return true
}

func (c StatsJobCommand) BuildResponse(responseLine string, body []byte) (CommandResponse, error) {
	switch {
	case strings.HasPrefix(responseLine, "OK"):
		return StatsJobCommandResponse{body}, nil

	case strings.EqualFold(responseLine, "NOT_FOUND"):
		return nil, ErrNotFound

	default:
		return nil, ErrUnexpectedResponse
	}
}
