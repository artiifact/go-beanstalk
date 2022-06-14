package beanstalk

import (
	"github.com/IvanLutokhin/go-beanstalk/internal/checker"
	"io"
	"net/textproto"
	"strconv"
	"strings"
)

var crnl = []byte{'\r', '\n'}

type Conn struct {
	conn    *textproto.Conn
	checker *checker.Checker
}

func NewConn(conn io.ReadWriteCloser) *Conn {
	return &Conn{
		conn:    textproto.NewConn(conn),
		checker: checker.NewChecker(conn),
	}
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

func (c *Conn) Check() error {
	return c.checker.Check()
}

func (c *Conn) WriteRequest(line string, body []byte) (uint, error) {
	id := c.conn.Next()

	c.conn.StartRequest(id)
	defer c.conn.EndRequest(id)

	if _, err := c.conn.W.Write([]byte(line)); err != nil {
		return 0, err
	}

	if _, err := c.conn.W.Write(crnl); err != nil {
		return 0, err
	}

	if body != nil {
		if _, err := c.conn.W.Write(body); err != nil {
			return 0, err
		}

		if _, err := c.conn.W.Write(crnl); err != nil {
			return 0, err
		}
	}

	if err := c.conn.W.Flush(); err != nil {
		return 0, err
	}

	return id, nil
}

func (c *Conn) ReadResponse(id uint, hasBody bool) (string, []byte, error) {
	c.conn.StartResponse(id)
	defer c.conn.EndResponse(id)

	line, err := c.conn.ReadLine()
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
		if _, err = io.ReadFull(c.conn.R, body); err != nil {
			return line, body, err
		}

		body = body[:n] // exclude CR NL
	}

	return line, body, nil
}
