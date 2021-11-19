package beanstalk

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type PutCommand struct {
	Priority uint32
	Delay    time.Duration
	TTR      time.Duration
	Data     []byte
}

type PutCommandResponse struct {
	ID int
}

func (c PutCommand) CommandLine() string {
	return fmt.Sprintf("put %d %0.f %0.f %d", c.Priority, c.Delay.Seconds(), c.TTR.Seconds(), len(c.Data))
}

func (c PutCommand) Body() []byte {
	return c.Data
}

func (c PutCommand) HasResponseBody() bool {
	return false
}

func (c PutCommand) BuildResponse(responseLine string, _ []byte) (CommandResponse, error) {
	switch {
	case strings.HasPrefix(responseLine, "INSERTED"):
		i := strings.LastIndex(responseLine, " ")
		if i == -1 {
			return nil, ErrUnexpectedResponse
		}

		id, err := strconv.Atoi(responseLine[i+1:])
		if err != nil {
			return nil, err
		}

		return PutCommandResponse{id}, nil

	case strings.HasPrefix(responseLine, "BURIED"):
		return nil, ErrBuried

	case strings.EqualFold(responseLine, "EXPECTED_CRLF"):
		return nil, ErrExpectedCRLF

	case strings.EqualFold(responseLine, "JOB_TOO_BIG"):
		return nil, ErrJobTooBig

	case strings.EqualFold(responseLine, "DRAINING"):
		return nil, ErrDraining

	default:
		return nil, ErrUnexpectedResponse
	}
}
