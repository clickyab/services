package rabbitmq

import "sync"

type publishCount struct {
	mu      sync.Mutex
	success int64
	failed  int64
}

//IncrSuccess increase success publish job
func (c *publishCount) IncrSuccess() {
	c.mu.Lock()
	c.success++
	c.mu.Unlock()
}

//IncrFailed increase failes publish job
func (c *publishCount) IncrFailed() {
	c.mu.Lock()
	c.failed++
	c.mu.Unlock()
}

//Value get publish message count
func (c *publishCount) Value() (int64, int64) {
	c.mu.Lock()
	suc := c.success
	fa := c.failed
	c.mu.Unlock()

	return suc, fa
}

type consumeCount struct {
	mu      sync.Mutex
	success int64
	failed  int64
}

//IncrSuccess increase success consume job
func (c *consumeCount) IncrSuccess() {
	c.mu.Lock()
	c.success++
	c.mu.Unlock()
}

//IncrFailed increase failes consume job
func (c *consumeCount) IncrFailed() {
	c.mu.Lock()
	c.failed++
	c.mu.Unlock()
}

//Value get counts of consumes job
func (c *consumeCount) Value() (int64, int64) {
	c.mu.Lock()
	suc := c.success
	fa := c.failed
	c.mu.Unlock()

	return suc, fa
}
