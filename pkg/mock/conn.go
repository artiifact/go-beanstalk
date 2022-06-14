package mock

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type ConnError struct {
	expected []byte
	actual   []byte
}

func (e ConnError) Error() string {
	return fmt.Sprintf("beanstalk: conn: expected %#v, got %#v", string(e.expected), string(e.actual))
}

type Conn struct {
	in  []string
	out []string
}

func NewConn(in, out []string) io.ReadWriteCloser {
	return &Conn{
		in:  in,
		out: out,
	}
}

func (c *Conn) Read(b []byte) (int, error) {
	if len(c.out) == 0 {
		return 0, io.EOF
	}

	n, err := strings.NewReader(c.out[0]).Read(b)
	if err != nil {
		return n, err
	}

	c.out = c.out[1:]

	return n, nil
}

func (c *Conn) Write(b []byte) (int, error) {
	if len(c.in) == 0 {
		return 0, io.EOF
	}

	expected := make([]byte, len(b))

	n, err := strings.NewReader(c.in[0]).Read(expected)
	if err != nil {
		return n, err
	}

	expected = expected[:n]

	if bytes.Compare(expected, b) != 0 {
		return 0, ConnError{expected: expected, actual: b}
	}

	c.in = c.in[1:]

	return n, nil
}

func (c *Conn) Close() error {
	switch {
	case len(c.in) > 0:
		return errors.New("beanstalk: conn: input buffer is not empty")

	case len(c.out) > 0:
		return errors.New("beanstalk: conn: output buffer is not empty")

	default:
		return nil
	}
}
