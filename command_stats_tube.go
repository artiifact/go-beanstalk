package beanstalk

import (
	"fmt"
	"strings"
)

type StatsTubeCommand struct {
	Tube string
}

type StatsTubeCommandResponse struct {
	Data []byte
}

func (c StatsTubeCommand) CommandLine() string {
	return fmt.Sprintf("stats-tube %s", c.Tube)
}

func (c StatsTubeCommand) Body() []byte {
	return nil
}

func (c StatsTubeCommand) HasResponseBody() bool {
	return true
}

func (c StatsTubeCommand) BuildResponse(responseLine string, body []byte) (CommandResponse, error) {
	switch {
	case strings.HasPrefix(responseLine, "OK"):
		return StatsTubeCommandResponse{body}, nil

	case strings.EqualFold(responseLine, "NOT_FOUND"):
		return nil, ErrNotFound

	default:
		return nil, ErrUnexpectedResponse
	}
}
