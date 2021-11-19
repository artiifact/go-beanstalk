package beanstalk

import (
	"strconv"
	"strings"
)

type ReserveCommand struct{}

type ReserveCommandResponse struct {
	ID   int
	Data []byte
}

func (c ReserveCommand) CommandLine() string {
	return "reserve"
}

func (c ReserveCommand) Body() []byte {
	return nil
}

func (c ReserveCommand) HasResponseBody() bool {
	return true
}

func (c ReserveCommand) BuildResponse(responseLine string, body []byte) (CommandResponse, error) {
	switch {
	case strings.HasPrefix(responseLine, "RESERVED"):
		fields := strings.Fields(responseLine)
		if len(fields) == 1 {
			return nil, ErrUnexpectedResponse
		}

		id, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}

		return ReserveCommandResponse{id, body}, nil

	case strings.EqualFold(responseLine, "DEADLINE_SOON"):
		return nil, ErrDeadlineSoon

	case strings.EqualFold(responseLine, "TIMED_OUT"):
		return nil, ErrTimedOut

	default:
		return nil, ErrUnexpectedResponse
	}
}
