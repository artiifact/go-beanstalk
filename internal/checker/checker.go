package checker

import "io"

type Checker struct {
	conn io.ReadWriteCloser
}

func NewChecker(conn io.ReadWriteCloser) *Checker {
	return &Checker{conn: conn}
}
