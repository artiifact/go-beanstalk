package beanstalk

import (
	"strconv"
	"strings"
)

type PeekDelayedCommand struct{}

type PeekDelayedCommandResponse struct {
	ID   int
	Data []byte
}

func (c PeekDelayedCommand) CommandLine() string {
	return "peek-delayed"
}

func (c PeekDelayedCommand) Body() []byte {
	return nil
}

func (c PeekDelayedCommand) HasResponseBody() bool {
	return true
}

func (c PeekDelayedCommand) BuildResponse(responseLine string, body []byte) (CommandResponse, error) {
	switch {
	case strings.HasPrefix(responseLine, "FOUND"):
		fields := strings.Fields(responseLine)
		if len(fields) == 1 {
			return nil, ErrUnexpectedResponse
		}

		id, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}

		return PeekDelayedCommandResponse{id, body}, nil

	case strings.EqualFold(responseLine, "NOT_FOUND"):
		return nil, ErrNotFound

	default:
		return nil, ErrUnexpectedResponse
	}
}
