package beanstalk

import "strings"

type StatsCommand struct{}

type StatsCommandResponse struct {
	Data []byte
}

func (c StatsCommand) CommandLine() string {
	return "stats"
}

func (c StatsCommand) Body() []byte {
	return nil
}

func (c StatsCommand) HasResponseBody() bool {
	return true
}

func (c StatsCommand) BuildResponse(responseLine string, body []byte) (CommandResponse, error) {
	switch {
	case strings.HasPrefix(responseLine, "OK"):
		return StatsCommandResponse{body}, nil

	default:
		return nil, ErrUnexpectedResponse
	}
}
