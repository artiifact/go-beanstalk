package beanstalk

import (
	"fmt"
	"strconv"
	"strings"
)

type PeekCommand struct {
	ID int
}

type PeekCommandResponse struct {
	ID   int
	Data []byte
}

func (c PeekCommand) CommandLine() string {
	return fmt.Sprintf("peek %d", c.ID)
}

func (c PeekCommand) Body() []byte {
	return nil
}

func (c PeekCommand) HasResponseBody() bool {
	return true
}

func (c PeekCommand) BuildResponse(responseLine string, body []byte) (CommandResponse, error) {
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

		return PeekCommandResponse{id, body}, nil

	case strings.EqualFold(responseLine, "NOT_FOUND"):
		return nil, ErrNotFound

	default:
		return nil, ErrUnexpectedResponse
	}
}
