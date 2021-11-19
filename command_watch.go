package beanstalk

import (
	"fmt"
	"strconv"
	"strings"
)

type WatchCommand struct {
	Tube string
}

type WatchCommandResponse struct {
	Count int
}

func (c WatchCommand) CommandLine() string {
	return fmt.Sprintf("watch %s", c.Tube)
}

func (c WatchCommand) Body() []byte {
	return nil
}

func (c WatchCommand) HasResponseBody() bool {
	return false
}

func (c WatchCommand) BuildResponse(responseLine string, _ []byte) (CommandResponse, error) {
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

		return WatchCommandResponse{count}, nil

	default:
		return nil, ErrUnexpectedResponse
	}
}
