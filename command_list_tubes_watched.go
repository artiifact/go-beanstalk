package beanstalk

import "strings"

type ListTubesWatchedCommand struct{}

type ListTubesWatchedCommandResponse struct {
	Data []byte
}

func (c ListTubesWatchedCommand) CommandLine() string {
	return "list-tubes-watched"
}

func (c ListTubesWatchedCommand) Body() []byte {
	return nil
}

func (c ListTubesWatchedCommand) HasResponseBody() bool {
	return true
}

func (c ListTubesWatchedCommand) BuildResponse(responseLine string, body []byte) (CommandResponse, error) {
	switch {
	case strings.HasPrefix(responseLine, "OK"):
		return ListTubesWatchedCommandResponse{body}, nil

	default:
		return nil, ErrUnexpectedResponse
	}
}
