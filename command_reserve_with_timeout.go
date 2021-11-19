package beanstalk

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ReserveWithTimeoutCommand struct {
	Timeout time.Duration
}

type ReserveWithTimeoutCommandResponse struct {
	ID   int
	Data []byte
}

func (c ReserveWithTimeoutCommand) CommandLine() string {
	return fmt.Sprintf("reserve-with-timeout %0.f", c.Timeout.Seconds())
}

func (c ReserveWithTimeoutCommand) Body() []byte {
	return nil
}

func (c ReserveWithTimeoutCommand) HasResponseBody() bool {
	return true
}

func (c ReserveWithTimeoutCommand) BuildResponse(responseLine string, body []byte) (CommandResponse, error) {
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

		return ReserveWithTimeoutCommandResponse{id, body}, nil

	case strings.EqualFold(responseLine, "DEADLINE_SOON"):
		return nil, ErrDeadlineSoon

	case strings.EqualFold(responseLine, "TIMED_OUT"):
		return nil, ErrTimedOut

	default:
		return nil, ErrUnexpectedResponse
	}
}
