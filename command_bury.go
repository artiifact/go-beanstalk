package beanstalk

import (
	"fmt"
	"strings"
)

type BuryCommand struct {
	ID       int
	Priority uint32
}

type BuryCommandResponse struct{}

func (c BuryCommand) CommandLine() string {
	return fmt.Sprintf("bury %d %d", c.ID, c.Priority)
}

func (c BuryCommand) Body() []byte {
	return nil
}

func (c BuryCommand) HasResponseBody() bool {
	return false
}

func (c BuryCommand) BuildResponse(responseLine string, _ []byte) (CommandResponse, error) {
	switch {
	case strings.EqualFold(responseLine, "BURIED"):
		return BuryCommandResponse{}, nil

	case strings.EqualFold(responseLine, "NOT_FOUND"):
		return nil, ErrNotFound

	default:
		return nil, ErrUnexpectedResponse
	}
}
