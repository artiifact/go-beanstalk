package beanstalk

import (
	"testing"
)

func newMockPool(min, max int) (*Pool, error) {
	return NewPool(min, max, func() (*Client, error) { return NewClient(NewMockConn(nil, nil)), nil })
}

func TestPool_Get(t *testing.T) {
	p, err := newMockPool(1, 2)
	if err != nil {
		t.Error(err)
	}

	if p.Len() != 1 {
		t.Fatal("expected 1, but got", p.Len())
	}

	// get exists client from pool
	client1, err := p.Get()
	if err != nil {
		t.Error(err)
	}

	if client1 == nil {
		t.Fatal("expected client, but got nil")
	}

	if p.Len() != 0 {
		t.Fatal("expected 0, but got", p.Len())
	}

	// use factory to create new client
	client2, err := p.Get()
	if err != nil {
		t.Error(err)
	}

	if client2 == nil {
		t.Fatal("expected client, but got nil")
	}

	if p.Len() != 0 {
		t.Fatal("expected 0, but got", p.Len())
	}
}

func TestPool_Put(t *testing.T) {
	p, err := newMockPool(0, 1)
	if err != nil {
		t.Error(err)
	}

	if p.Len() != 0 {
		t.Fatal("expected 0, but got", p.Len())
	}

	// use factory to create new client
	client1, err := p.Get()
	if err != nil {
		t.Error(err)
	}

	if client1 == nil {
		t.Fatal("expected client, but got nil")
	}

	if p.Len() != 0 {
		t.Fatal("expected 0, but got", p.Len())
	}

	// use factory to create new client
	client2, err := p.Get()
	if err != nil {
		t.Error(err)
	}

	if client2 == nil {
		t.Fatal("expected client, but got nil")
	}

	if p.Len() != 0 {
		t.Fatal("expected 0, but got", p.Len())
	}

	// put to pool
	if err = p.Put(client1); err != nil {
		t.Error(err)
	}

	if p.Len() != 1 {
		t.Fatal("expected 1, but got", p.Len())
	}

	// put to pool and close client (max limit)
	if err = p.Put(client2); err != nil {
		t.Error(err)
	}

	if p.Len() != 1 {
		t.Fatal("expected 1, but got", p.Len())
	}
}
