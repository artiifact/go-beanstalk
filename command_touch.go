package beanstalk

import (
	"fmt"
	"strings"
)

type TouchCommand struct {
	ID int
}

type TouchCommandResponse struct{}

func (c TouchCommand) CommandLine() string {
	return fmt.Sprintf("touch %d", c.ID)
}

func (c TouchCommand) Body() []byte {
	return nil
}

func (c TouchCommand) HasResponseBody() bool {
	return false
}

func (c TouchCommand) BuildResponse(responseLine string, _ []byte) (CommandResponse, error) {
	switch {
	case strings.EqualFold(responseLine, "TOUCHED"):
		return TouchCommandResponse{}, nil

	case strings.EqualFold(responseLine, "NOT_FOUND"):
		return nil, ErrNotFound

	default:
		return nil, ErrUnexpectedResponse
	}
}
