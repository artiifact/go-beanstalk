//go:build windows
// +build windows

package checker

import (
	"errors"
	"io"
	"syscall"
)

func (c *Checker) Check() error {
	sysConn, ok := c.conn.(syscall.Conn)
	if !ok {
		return nil
	}

	rawConn, err := sysConn.SyscallConn()
	if err != nil {
		return err
	}

	var sysErr error
	err = rawConn.Read(func(fd uintptr) bool {
		var buffer [1]byte

		n, err := syscall.Read(syscall.Handle(fd), buffer[:])
		switch {
		case n == 0 && err == nil:
			sysErr = io.EOF
		case n > 0:
			sysErr = errors.New("unexpected read from socket")
		case err == syscall.EAGAIN || err == syscall.EWOULDBLOCK:
			sysErr = nil
		default:
			sysErr = err
		}

		return true
	})
	if err != nil {
		return err
	}

	return sysErr
}
