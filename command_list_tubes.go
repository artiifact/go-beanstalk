package beanstalk

import "strings"

type ListTubesCommand struct{}

type ListTubesCommandResponse struct {
	Data []byte
}

func (c ListTubesCommand) CommandLine() string {
	return "list-tubes"
}

func (c ListTubesCommand) Body() []byte {
	return nil
}

func (c ListTubesCommand) HasResponseBody() bool {
	return true
}

func (c ListTubesCommand) BuildResponse(responseLine string, body []byte) (CommandResponse, error) {
	switch {
	case strings.HasPrefix(responseLine, "OK"):
		return ListTubesCommandResponse{body}, nil

	default:
		return nil, ErrUnexpectedResponse
	}
}
