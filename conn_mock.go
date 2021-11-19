package beanstalk

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type MockConnError struct {
	expected []byte
	actual   []byte
}

func (e MockConnError) Error() string {
	return fmt.Sprintf("beanstalk: mock conn: expected %#v, got %#v", string(e.expected), string(e.actual))
}

type MockConn struct {
	in  []string
	out []string
}

func NewMockConn(in, out []string) io.ReadWriteCloser {
	return &MockConn{
		in:  in,
		out: out,
	}
}

func (c *MockConn) Read(b []byte) (int, error) {
	if len(c.out) == 0 {
		return 0, errors.New("beanstalk: mock conn: EOF on read")
	}

	n, err := strings.NewReader(c.out[0]).Read(b)
	if err != nil {
		return n, err
	}

	c.out = c.out[1:]

	return n, nil
}

func (c *MockConn) Write(b []byte) (int, error) {
	if len(c.in) == 0 {
		return 0, errors.New("beanstalk: mock conn: EOF on write")
	}

	expected := make([]byte, len(b))

	n, err := strings.NewReader(c.in[0]).Read(expected)
	if err != nil {
		return n, err
	}

	expected = expected[:n]

	if bytes.Compare(expected, b) != 0 {
		return 0, MockConnError{expected, b}
	}

	c.in = c.in[1:]

	return n, nil
}

func (c MockConn) Close() error {
	switch {
	case len(c.in) > 0:
		return errors.New("beanstalk: mock conn: not empty input buffer")

	case len(c.out) > 0:
		return errors.New("beanstalk: mock conn: not empty output buffer")

	default:
		return nil
	}
}
