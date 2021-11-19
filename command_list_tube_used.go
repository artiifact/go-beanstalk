package beanstalk

import (
	"strings"
)

type ListTubeUsedCommand struct{}

type ListTubeUsedCommandResponse struct {
	Tube string
}

func (c ListTubeUsedCommand) CommandLine() string {
	return "list-tube-used"
}

func (c ListTubeUsedCommand) Body() []byte {
	return nil
}

func (c ListTubeUsedCommand) HasResponseBody() bool {
	return false
}

func (c ListTubeUsedCommand) BuildResponse(responseLine string, _ []byte) (CommandResponse, error) {
	switch {
	case strings.HasPrefix(responseLine, "USING"):
		i := strings.LastIndex(responseLine, " ")
		if i == -1 {
			return nil, ErrUnexpectedResponse
		}

		return ListTubeUsedCommandResponse{responseLine[i+1:]}, nil

	default:
		return nil, ErrUnexpectedResponse
	}
}
