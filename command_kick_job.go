package beanstalk

import (
	"fmt"
	"strings"
)

type KickJobCommand struct {
	ID int
}

type KickJobCommandResponse struct{}

func (c KickJobCommand) CommandLine() string {
	return fmt.Sprintf("kick-job %d", c.ID)
}

func (c KickJobCommand) Body() []byte {
	return nil
}

func (c KickJobCommand) HasResponseBody() bool {
	return false
}

func (c KickJobCommand) BuildResponse(responseLine string, _ []byte) (CommandResponse, error) {
	switch {
	case strings.EqualFold(responseLine, "KICKED"):
		return KickJobCommandResponse{}, nil

	case strings.EqualFold(responseLine, "NOT_FOUND"):
		return nil, ErrNotFound

	default:
		return nil, ErrUnexpectedResponse
	}
}
