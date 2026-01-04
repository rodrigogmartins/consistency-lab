package replication

import (
	"math/rand"
	"sync"
	"time"
)

type Chaos struct {
	DropRate float64 // 0..1
	MinDelay time.Duration
	MaxDelay time.Duration
	Rand     *rand.Rand
	mu       sync.Mutex
}

func (c *Chaos) ShouldDrop() bool {
	if c.DropRate <= 0 {
		return false
	}
	c.mu.Lock()
	v := c.Rand.Float64()
	c.mu.Unlock()
	return v < c.DropRate
}

func (c *Chaos) Delay() time.Duration {
	if c.MaxDelay <= 0 || c.MaxDelay <= c.MinDelay {
		return c.MinDelay
	}
	delta := c.MaxDelay - c.MinDelay

	c.mu.Lock()
	n := time.Duration(c.Rand.Int63n(int64(delta) + 1))
	c.mu.Unlock()

	return c.MinDelay + n
}
