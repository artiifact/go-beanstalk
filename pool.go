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
	ErrIllegalClient      = errors.New("beanstalk: pool: illegal client")
)

type Pool interface {
	Open(ctx context.Context) error
	Close(ctx context.Context) error
	Get() (Client, error)
	Put(client Client) error
	Len() int
}

type PoolOptions struct {
	Dialer      func() (*DefaultClient, error)
	Logger      Logger
	Capacity    int
	MaxAge      time.Duration
	IdleTimeout time.Duration
}

type DefaultPool struct {
	options   *PoolOptions
	clients   []*DefaultClient
	triggerCh chan struct{}
	closeCh   chan struct{}
	closed    int32
	mutex     sync.RWMutex
}

func NewDefaultPool(options *PoolOptions) *DefaultPool {
	if options.Logger == nil {
		options.Logger = NopLogger
	}

	if options.Capacity < 1 {
		options.Capacity = 1
	}

	if options.MaxAge < 0 {
		options.MaxAge = 0
	}

	if options.IdleTimeout < 0 {
		options.IdleTimeout = 0
	}

	return &DefaultPool{
		options:   options,
		clients:   make([]*DefaultClient, 0, options.Capacity),
		triggerCh: make(chan struct{}),
		closeCh:   make(chan struct{}),
		closed:    1,
	}
}

func (p *DefaultPool) Open(ctx context.Context) error {
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

func (p *DefaultPool) Close(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&p.closed, 0, 1) {
		return ErrClosedPool
	}

	p.closeCh <- struct{}{}

	doneCh := make(chan struct{})

	go func() {
		defer close(doneCh)

		p.mutex.Lock()
		defer p.mutex.Unlock()

		for _, client := range p.clients {
			if err := client.Close(); err != nil {
				p.options.Logger.Log(ErrorLogLevel, "Failed to close client", map[string]interface{}{"error": err})
			}
		}

		p.clients = p.clients[:0]
	}()

	select {
	case <-doneCh:
		p.options.Logger.Log(InfoLogLevel, "Pool was closed", nil)

		return nil
	case <-ctx.Done():
		return fmt.Errorf("beanstalk: pool: %v", ctx.Err())
	}
}

func (p *DefaultPool) Get() (Client, error) {
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

func (p *DefaultPool) Put(client Client) error {
	if p.isClosed() {
		return ErrClosedPool
	}

	p.options.Logger.Log(DebugLogLevel, "Tries to return client", nil)

	defaultClient, ok := client.(*DefaultClient)
	if !ok {
		p.options.Logger.Log(DebugLogLevel, "Failed to return client", map[string]interface{}{"error": ErrIllegalClient})

		return ErrIllegalClient
	}

	if p.Len() >= p.options.Capacity || !p.checkClient(defaultClient) {
		p.options.Logger.Log(DebugLogLevel, "Closes stale client", nil)

		if err := defaultClient.Close(); err != nil {
			p.options.Logger.Log(ErrorLogLevel, "Failed to close client", map[string]interface{}{"error": err})
		}

		return nil
	}

	p.mutex.Lock()
	p.clients = append(p.clients, defaultClient)
	p.mutex.Unlock()

	p.options.Logger.Log(DebugLogLevel, "Client was returned", nil)

	return nil
}

func (p *DefaultPool) Len() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return len(p.clients)
}

func (p *DefaultPool) isClosed() bool {
	return atomic.LoadInt32(&p.closed) == 1
}

func (p *DefaultPool) refillClients() {
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

func (p *DefaultPool) createClient() (*DefaultClient, error) {
	if p.options.Dialer == nil {
		return nil, ErrDialerNotSpecified
	}

	return p.options.Dialer()
}

func (p *DefaultPool) createAndPutClient() {
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

func (p *DefaultPool) checkClient(client *DefaultClient) bool {
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
