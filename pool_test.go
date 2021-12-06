package beanstalk

import (
	"testing"
)

func newMockPool(capacity int, open bool) (Pool, error) {
	return NewPool(func() (Client, error) { return NewClient(NewMockConn(nil, nil)), nil }, capacity, open)
}

func TestPool_Get(t *testing.T) {
	p, err := newMockPool(1, true)
	if err != nil {
		t.Error(err)
	}

	if p.Len() != 1 {
		t.Fatal("expected 1, but got", p.Len())
	}

	if err = p.Open(); err == nil {
		t.Fatal("expected error, but got nil")
	}

	// get exists client from pool
	client1, err := p.Get()
	if err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
	}

	if client2 == nil {
		t.Fatal("expected client, but got nil")
	}

	if p.Len() != 0 {
		t.Fatal("expected 0, but got", p.Len())
	}

	if err = p.Close(); err != nil {
		t.Fatal(err)
	}

	if _, err = p.Get(); err == nil {
		t.Fatal("expected error, but got nil")
	}
}

func TestPool_Put(t *testing.T) {
	p, err := newMockPool(1, false)
	if err != nil {
		t.Error(err)
	}

	if p.Len() != 0 {
		t.Fatal("expected 0, but got", p.Len())
	}

	if err = p.Open(); err != nil {
		t.Fatal(err)
	}

	// get exists client from pool
	client1, err := p.Get()
	if err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
	}

	if client2 == nil {
		t.Fatal("expected client, but got nil")
	}

	if p.Len() != 0 {
		t.Fatal("expected 0, but got", p.Len())
	}

	// put to pool
	if err = p.Put(client1); err != nil {
		t.Fatal(err)
	}

	if p.Len() != 1 {
		t.Fatal("expected 1, but got", p.Len())
	}

	// put to pool and close client (max limit)
	if err = p.Put(client2); err != nil {
		t.Fatal(err)
	}

	if p.Len() != 1 {
		t.Fatal("expected 1, but got", p.Len())
	}

	if err = p.Close(); err != nil {
		t.Fatal(err)
	}

	if err = p.Put(nil); err == nil {
		t.Fatal("expected error, but got nil")
	}
}
