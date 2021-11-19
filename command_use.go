package beanstalk

import (
	"fmt"
	"strings"
)

type UseCommand struct {
	Tube string
}

type UseCommandResponse struct {
	Tube string
}

func (c UseCommand) CommandLine() string {
	return fmt.Sprintf("use %s", c.Tube)
}

func (c UseCommand) Body() []byte {
	return nil
}

func (c UseCommand) HasResponseBody() bool {
	return false
}

func (c UseCommand) BuildResponse(responseLine string, _ []byte) (CommandResponse, error) {
	switch {
	case strings.HasPrefix(responseLine, "USING"):
		i := strings.LastIndex(responseLine, " ")
		if i == -1 {
			return nil, ErrUnexpectedResponse
		}

		return UseCommandResponse{responseLine[i+1:]}, nil

	default:
		return nil, ErrUnexpectedResponse
	}
}
