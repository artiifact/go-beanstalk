package beanstalk

import (
	"fmt"
	"strconv"
	"strings"
)

type KickCommand struct {
	Bound int
}

type KickCommandResponse struct {
	Count int
}

func (c KickCommand) CommandLine() string {
	return fmt.Sprintf("kick %d", c.Bound)
}

func (c KickCommand) Body() []byte {
	return nil
}

func (c KickCommand) HasResponseBody() bool {
	return false
}

func (c KickCommand) BuildResponse(responseLine string, _ []byte) (CommandResponse, error) {
	switch {
	case strings.HasPrefix(responseLine, "KICKED"):
		i := strings.LastIndex(responseLine, " ")
		if i == -1 {
			return nil, ErrUnexpectedResponse
		}

		count, err := strconv.Atoi(responseLine[i+1:])
		if err != nil {
			return nil, err
		}

		return KickCommandResponse{count}, nil

	default:
		return nil, ErrUnexpectedResponse
	}
}
