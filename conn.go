package beanstalk

import (
	"io"
	"net/textproto"
	"strconv"
	"strings"
)

var crnl = []byte{'\r', '\n'}

type Conn struct {
	Text *textproto.Conn
}

func NewConn(conn *textproto.Conn) *Conn {
	return &Conn{conn}
}

func (c *Conn) Close() error {
	return c.Text.Close()
}

func (c *Conn) WriteRequest(line string, body []byte) (uint, error) {
	id := c.Text.Next()

	c.Text.StartRequest(id)
	defer c.Text.EndRequest(id)

	if _, err := c.Text.W.Write([]byte(line)); err != nil {
		return 0, err
	}

	if _, err := c.Text.W.Write(crnl); err != nil {
		return 0, err
	}

	if body != nil {
		if _, err := c.Text.W.Write(body); err != nil {
			return 0, err
		}

		if _, err := c.Text.W.Write(crnl); err != nil {
			return 0, err
		}
	}

	if err := c.Text.W.Flush(); err != nil {
		return 0, err
	}

	return id, nil
}

func (c *Conn) ReadResponse(id uint, hasBody bool) (string, []byte, error) {
	c.Text.StartResponse(id)
	defer c.Text.EndResponse(id)

	line, err := c.Text.ReadLine()
	if err != nil {
		return line, nil, err
	}

	var body []byte
	if hasBody {
		i := strings.LastIndex(line, " ")
		if i == -1 {
			return line, body, nil
		}

		n, err := strconv.Atoi(line[i+1:])
		if err != nil {
			return line, body, err
		}

		body = make([]byte, n+2) // include CR NL
		if _, err = io.ReadFull(c.Text.R, body); err != nil {
			return line, body, err
		}

		body = body[:n] // exclude CR NL
	}

	return line, body, nil
}
