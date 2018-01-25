package main

import (
	"io"
	"sync"

	"log"

	"github.com/pkg/errors"
)

type Pool struct {
	m         sync.Mutex
	resources chan io.Closer
	factory   func() (io.Closer, error)
	closed    bool
}

var ErrPoolClosed = errors.New("Pool is closed.")

func NewPool(f func() (io.Closer, error), count int) (*Pool, error) {
	if count <= 0 {
		return nil, errors.New("Pool is too small")
	}

	return &Pool{
		resources: make(chan io.Closer, count),
		factory:   f,
	}, nil
}

func (p *Pool) Acquire() (io.Closer, error) {
	select {
	case r, ok := <-p.resources:
		if !ok {
			return nil, ErrPoolClosed
		}
		return r, nil
	default:
		return p.factory()
	}
}

func (p *Pool) Release(r io.Closer) {
	p.m.Lock()
	defer p.m.Unlock()

	if p.closed {
		r.Close()
		return
	}

	select {
	case p.resources <- r:
		log.Println("Placed a resource in the queue")
	default:
		log.Println("Closing resource - queue is already of capacity")
		r.Close()
	}
}

func (p *Pool) Close() {
	p.m.Lock()
	defer p.m.Unlock()

	if p.closed {
		return
	}

	p.closed = true

	close(p.resources)

	for r := range p.resources {
		r.Close()
	}
}
