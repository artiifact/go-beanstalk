package beanstalk

import (
	"fmt"
	"strings"
)

type DeleteCommand struct {
	ID int
}

type DeleteCommandResponse struct{}

func (c DeleteCommand) CommandLine() string {
	return fmt.Sprintf("delete %d", c.ID)
}

func (c DeleteCommand) Body() []byte {
	return nil
}

func (c DeleteCommand) HasResponseBody() bool {
	return false
}

func (c DeleteCommand) BuildResponse(responseLine string, _ []byte) (CommandResponse, error) {
	switch {
	case strings.EqualFold(responseLine, "DELETED"):
		return DeleteCommandResponse{}, nil

	case strings.EqualFold(responseLine, "NOT_FOUND"):
		return nil, ErrNotFound

	default:
		return nil, ErrUnexpectedResponse
	}
}
