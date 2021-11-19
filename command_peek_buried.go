package beanstalk

import (
	"strconv"
	"strings"
)

type PeekBuriedCommand struct{}

type PeekBuriedCommandResponse struct {
	ID   int
	Data []byte
}

func (c PeekBuriedCommand) CommandLine() string {
	return "peek-buried"
}

func (c PeekBuriedCommand) Body() []byte {
	return nil
}

func (c PeekBuriedCommand) HasResponseBody() bool {
	return true
}

func (c PeekBuriedCommand) BuildResponse(responseLine string, body []byte) (CommandResponse, error) {
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

		return PeekBuriedCommandResponse{id, body}, nil

	case strings.EqualFold(responseLine, "NOT_FOUND"):
		return nil, ErrNotFound

	default:
		return nil, ErrUnexpectedResponse
	}
}
