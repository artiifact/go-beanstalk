package beanstalk

import "errors"

var (
	ErrBadFormat          = errors.New("beanstalk: bad format")
	ErrBuried             = errors.New("beanstalk: buried")
	ErrDeadlineSoon       = errors.New("beanstalk: deadline soon")
	ErrDraining           = errors.New("beanstalk: draining")
	ErrExpectedCRLF       = errors.New("beanstalk: expected CRLF")
	ErrInternalError      = errors.New("beanstalk: internal error")
	ErrJobTooBig          = errors.New("beanstalk: job too big")
	ErrNotFound           = errors.New("beanstalk: not found")
	ErrNotIgnored         = errors.New("beanstalk: not ignored")
	ErrOutOfMemory        = errors.New("beanstalk: out of memory")
	ErrTimedOut           = errors.New("beanstalk: timed out")
	ErrUnknownCommand     = errors.New("beanstalk: unknown command")
	ErrMalformedCommand   = errors.New("beanstalk: malformed command")
	ErrUnexpectedResponse = errors.New("beanstalk: unexpected response")
)
