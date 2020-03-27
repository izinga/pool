package pool

import (
	"sync"
)

// PoolConn is a wrapper around Client to modify the the behavior of
// Client's Close() method.
type PoolConn struct {
	Client
	mu       sync.RWMutex
	c        *channelPool
	unusable bool
}

// Close() puts the given connects back to the pool instead of closing it.
func (p *PoolConn) Close() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.unusable {
		if p.Client != nil {
			return p.Client.Close()
		}
		return nil
	}
	return p.c.put(p.Client)
}

// MarkUnusable() marks the connection not usable any more, to let the pool close it instead of returning it to pool.
func (p *PoolConn) MarkUnusable() {
	p.mu.Lock()
	p.unusable = true
	p.mu.Unlock()
}

// newConn wraps a standard Client to a poolConn Client.
func (c *channelPool) wrapConn(conn Client) Client {
	p := &PoolConn{c: c}
	p.Client = conn
	return p
}
