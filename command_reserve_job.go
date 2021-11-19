package beanstalk

import (
	"fmt"
	"strconv"
	"strings"
)

type ReserveJobCommand struct {
	ID int
}

type ReserveJobCommandResponse struct {
	ID   int
	Data []byte
}

func (c ReserveJobCommand) CommandLine() string {
	return fmt.Sprintf("reserve-job %d", c.ID)
}

func (c ReserveJobCommand) Body() []byte {
	return nil
}

func (c ReserveJobCommand) HasResponseBody() bool {
	return true
}

func (c ReserveJobCommand) BuildResponse(responseLine string, body []byte) (CommandResponse, error) {
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

		return ReserveJobCommandResponse{id, body}, nil

	case strings.EqualFold(responseLine, "NOT_FOUND"):
		return nil, ErrNotFound

	default:
		return nil, ErrUnexpectedResponse
	}
}
