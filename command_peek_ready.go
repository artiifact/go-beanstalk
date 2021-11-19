package beanstalk

import (
	"strconv"
	"strings"
)

type PeekReadyCommand struct{}

type PeekReadyCommandResponse struct {
	ID   int
	Data []byte
}

func (c PeekReadyCommand) CommandLine() string {
	return "peek-ready"
}

func (c PeekReadyCommand) Body() []byte {
	return nil
}

func (c PeekReadyCommand) HasResponseBody() bool {
	return true
}

func (c PeekReadyCommand) BuildResponse(responseLine string, body []byte) (CommandResponse, error) {
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

		return PeekReadyCommandResponse{id, body}, nil

	case strings.EqualFold(responseLine, "NOT_FOUND"):
		return nil, ErrNotFound

	default:
		return nil, ErrUnexpectedResponse
	}
}
