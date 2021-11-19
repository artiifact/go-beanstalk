package beanstalk

type Command interface {
	CommandLine() string
	Body() []byte
	HasResponseBody() bool
}

type CommandResponse interface{}

type CommandResponseBuilder interface {
	BuildResponse(responseLine string, data []byte) (CommandResponse, error)
}
