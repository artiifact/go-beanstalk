package beanstalk

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrAlreadyOpenedPool  = errors.New("beanstalk: pool: already opened")
	ErrClosedPool         = errors.New("beanstalk: pool: closed")
	ErrDialerNotSpecified = errors.New("beanstalk: pool: dialer not specified")
)

type PoolOptions struct {
	Dialer      func() (*Client, error)
	Logger      Logger
	Capacity    int
	MaxAge      time.Duration
	IdleTimeout time.Duration
}

type Pool struct {
	options   *PoolOptions
	clients   []*Client
	triggerCh chan struct{}
	closeCh   chan struct{}
	closed    int32
	mutex     sync.RWMutex
}

func NewPool(options *PoolOptions) *Pool {
	if options.Logger == nil {
		options.Logger = NopLogger
	}

	if options.Capacity < 1 {
		options.Capacity = 1
	}

	return &Pool{
		options:   options,
		clients:   make([]*Client, 0, options.Capacity),
		triggerCh: make(chan struct{}),
		closeCh:   make(chan struct{}),
		closed:    1,
	}
}

func (p *Pool) Open(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&p.closed, 1, 0) {
		return ErrAlreadyOpenedPool
	}

	var wg sync.WaitGroup

	for i := 0; i < p.options.Capacity; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			p.createAndPutClient()
		}()
	}

	doneCh := make(chan struct{})

	go func() {
		defer close(doneCh)

		wg.Wait()
	}()

	select {
	case <-doneCh:
		go p.refillClients()

		p.options.Logger.Log(InfoLogLevel, "Pool was opened", nil)

		return nil
	case <-ctx.Done():
		return fmt.Errorf("beanstalk: pool: %v", ctx.Err())
	}
}

func (p *Pool) Close(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&p.closed, 0, 1) {
		return ErrClosedPool
	}

	p.closeCh <- struct{}{}

	doneCh := make(chan struct{})

	go func() {
		defer close(doneCh)

		p.mutex.Lock()
		for _, client := range p.clients {
			if err := client.Close(); err != nil {
				p.options.Logger.Log(ErrorLogLevel, "Failed to close client", map[string]interface{}{"error": err})
			}
		}

		p.clients = p.clients[:0]
		p.mutex.Unlock()
	}()

	select {
	case <-doneCh:
		p.options.Logger.Log(InfoLogLevel, "Pool was closed", nil)

		return nil
	case <-ctx.Done():
		return fmt.Errorf("beanstalk: pool: %v", ctx.Err())
	}
}

func (p *Pool) Get() (*Client, error) {
	if p.isClosed() {
		return nil, ErrClosedPool
	}

	for {
		if p.Len() == 0 {
			break
		}

		p.options.Logger.Log(DebugLogLevel, "Tries to fetch client", nil)

		p.mutex.Lock()
		client := p.clients[0]
		p.clients = append(p.clients[:0], p.clients[1:]...)
		p.mutex.Unlock()

		if !p.checkClient(client) {
			p.options.Logger.Log(DebugLogLevel, "Closes stale client", nil)

			if err := client.Close(); err != nil {
				p.options.Logger.Log(ErrorLogLevel, "Failed to close client", map[string]interface{}{"error": err})
			}

			continue
		}

		p.options.Logger.Log(DebugLogLevel, "Client was fetched", nil)

		return client, nil
	}

	p.triggerCh <- struct{}{}

	p.options.Logger.Log(DebugLogLevel, "Gets client by factory method", nil)

	return p.createClient()
}

func (p *Pool) Put(client *Client) error {
	if p.isClosed() {
		return ErrClosedPool
	}

	p.options.Logger.Log(DebugLogLevel, "Tries to return client", nil)

	if p.Len() >= p.options.Capacity || !p.checkClient(client) {
		p.options.Logger.Log(DebugLogLevel, "Closes stale client", nil)

		if err := client.Close(); err != nil {
			p.options.Logger.Log(ErrorLogLevel, "Failed to close client", map[string]interface{}{"error": err})
		}

		return nil
	}

	p.mutex.Lock()
	p.clients = append(p.clients, client)
	p.mutex.Unlock()

	p.options.Logger.Log(DebugLogLevel, "Client was returned", nil)

	return nil
}

func (p *Pool) Len() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return len(p.clients)
}

func (p *Pool) isClosed() bool {
	return atomic.LoadInt32(&p.closed) == 1
}

func (p *Pool) refillClients() {
	ticker := time.NewTicker(60 * time.Second)

	defer ticker.Stop()

	for {
		select {
		case <-p.closeCh:
			return

		case <-p.triggerCh:
		case <-ticker.C:
		}

		var wg sync.WaitGroup

		for i := p.Len(); i < p.options.Capacity; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				p.createAndPutClient()
			}()
		}

		wg.Wait()
	}
}

func (p *Pool) createClient() (*Client, error) {
	if p.options.Dialer == nil {
		return nil, ErrDialerNotSpecified
	}

	return p.options.Dialer()
}

func (p *Pool) createAndPutClient() {
	client, err := p.createClient()
	if err != nil {
		p.options.Logger.Log(ErrorLogLevel, "Failed to create client", map[string]interface{}{"error": err})

		return
	}

	p.options.Logger.Log(DebugLogLevel, "Client was created", nil)

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.isClosed() || len(p.clients) >= p.options.Capacity {
		p.options.Logger.Log(DebugLogLevel, "Closes stale client", nil)

		if err = client.Close(); err != nil {
			p.options.Logger.Log(ErrorLogLevel, "Failed to close client", map[string]interface{}{"error": err})
		}

		return
	}

	p.clients = append(p.clients, client)
}

func (p *Pool) checkClient(client *Client) bool {
	now := time.Now()

	if p.options.MaxAge > 0 && now.Sub(client.CreatedAt()) >= p.options.MaxAge {
		return false
	}

	if p.options.IdleTimeout > 0 && now.Sub(client.UsedAt()) >= p.options.IdleTimeout {
		return false
	}

	if client.ClosedAt().Unix() > 0 {
		return false
	}

	if err := client.Check(); err != nil {
		return false
	}

	return true
}
