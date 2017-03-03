package socket

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	// ErrClosed is the error resulting if the pool is closed via pool.Close().
	ErrClosed = errors.New("pool is closed")
	// ErrMaxedOut is the pool is already maxed out
	ErrMaxedOut = errors.New("the pool is already maxed out")
	// ErrCapacitySettings means the capacity arrguments is invalid
	ErrCapacitySettings = errors.New("invalid capacity settings")
)

// NewFactory ...
func NewFactory(addres string) Factory {
	return func() (conn net.Conn, err error) {
		return net.DialTimeout("tcp", addres, time.Minute)
	}
}

// Factory is a function to create new connections.
type Factory func() (net.Conn, error)

// NewPool create a new socket pool with an initial capacity and maximum capacity
func NewPool(initialCap, maxCap int, factory Factory) (pool *Pool, err error) {
	if initialCap < 0 || maxCap <= 0 || initialCap > maxCap {
		return nil, errors.New("invalid capacity settings")
	}
	pool = &Pool{
		sockets: make([]*Socket, 0),
		maxCap:  maxCap,
		factory: factory,
	}
	// create initial connections, if something goes wrong, just close the pool error out.
	for i := 0; i < initialCap; i++ {
		conn, err := factory()
		if err != nil || conn == nil {
			pool.Close()
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		pool.sockets = append(pool.sockets, pool.wrapConn(conn))
	}
	return
}

// Pool ...
type Pool struct {
	mu      sync.Mutex
	sockets []*Socket
	factory Factory
	maxCap  int
}

// Get returns a new connection from the pool.
func (p *Pool) Get() (socket *Socket, err error) {
	p.mu.Lock()
	if p.sockets == nil {
		return nil, ErrClosed
	}
	defer p.mu.Unlock()
	for _, socket = range p.sockets {
		if socket != nil && !socket.IsUse() {
			socket.Use()
			return
		}
	}
	if p.Len() >= p.maxCap {
		return nil, ErrMaxedOut
	}
	conn, err := p.factory()
	if err != nil || conn == nil {
		for _, conn := range p.sockets {
			conn.Close()
		}
		p.sockets = nil
		return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
	}
	socket = p.wrapConn(conn)
	socket.Use()
	p.sockets = append(p.sockets, socket)
	return
}

// Cap returns the maximum capacity of the pool
func (p *Pool) Cap() int {
	return p.maxCap
}

// Len returns the current capacity of the pool.
func (p *Pool) Len() int {
	return len(p.sockets)
}

// Close ...
func (p *Pool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, conn := range p.sockets {
		conn.Close()
	}
	p.sockets = nil
}

func (p *Pool) wrapConn(conn net.Conn) *Socket {
	return &Socket{Conn: conn}
}
