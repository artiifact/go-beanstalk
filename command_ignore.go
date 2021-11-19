package beanstalk

import (
	"fmt"
	"strconv"
	"strings"
)

type IgnoreCommand struct {
	Tube string
}

type IgnoreCommandResponse struct {
	Count int
}

func (c IgnoreCommand) CommandLine() string {
	return fmt.Sprintf("ignore %s", c.Tube)
}

func (c IgnoreCommand) Body() []byte {
	return nil
}

func (c IgnoreCommand) HasResponseBody() bool {
	return false
}

func (c IgnoreCommand) BuildResponse(responseLine string, _ []byte) (CommandResponse, error) {
	switch {
	case strings.HasPrefix(responseLine, "WATCHING"):
		i := strings.LastIndex(responseLine, " ")
		if i == -1 {
			return nil, ErrUnexpectedResponse
		}

		count, err := strconv.Atoi(responseLine[i+1:])
		if err != nil {
			return nil, err
		}

		return IgnoreCommandResponse{count}, nil

	case strings.EqualFold(responseLine, "NOT_IGNORED"):
		return nil, ErrNotIgnored

	default:
		return nil, ErrUnexpectedResponse
	}
}
