package beanstalk

import "errors"

type Pool struct {
	newClient func() (*Client, error)
	clients   chan *Client
}

func NewPool(min, max int, newClient func() (*Client, error)) (*Pool, error) {
	if max == 0 || min > max {
		return nil, errors.New("beanstalk: pool: invalid capacity")
	}

	if newClient == nil {
		return nil, errors.New("beanstalk: pool: factory function is nil")
	}

	p := &Pool{
		newClient: newClient,
		clients:   make(chan *Client, max),
	}

	for i := 0; i < min; i++ {
		client, err := p.newClient()
		if err != nil {
			return nil, err
		}

		p.clients <- client
	}

	return p, nil
}

func NewDefaultPool(address string, min, max int) (*Pool, error) {
	return NewPool(min, max, func() (*Client, error) { return Dial(address) })
}

func (p *Pool) Close() error {
	close(p.clients)

	for client := range p.clients {
		client.Close()
	}

	return nil
}

func (p *Pool) Get() (*Client, error) {
	select {
	case client := <-p.clients:
		return client, nil

	default:
		return p.newClient()
	}
}

func (p *Pool) Put(client *Client) error {
	select {
	case p.clients <- client:
		return nil

	default:
		return client.Close()
	}
}

func (p *Pool) Len() int {
	return len(p.clients)
}
