package mock

import (
	"sync"
	"time"

	"github.com/clickyab/services/kv"
)

var (
	all = make(map[string]chan struct{})
)

type locker struct {
	key   string
	ttl   time.Duration
	state *sync.Mutex
}

func (l *locker) Lock() {
	lock.Lock()
	this, ok := all[l.key]
	if !ok {
		this = make(chan struct{}, 1)
		all[l.key] = this
	}
	lock.Unlock()
	go func() {
		<-time.After(l.ttl)
		l.Unlock()
	}()
	this <- struct{}{}
}

func (l *locker) Unlock() {
	lock.Lock()
	this, ok := all[l.key]
	if !ok {
		this = make(chan struct{}, 1)
		all[l.key] = this
	}
	lock.Unlock()
	select {
	case <-this:
	default:
	}
}

func (l locker) Resource() string {
	return l.key
}

func (l locker) TTL() time.Duration {
	return l.ttl
}

// NewMockDistributedLocker return a new locker, its local not distributed
func NewMockDistributedLocker(resource string, d time.Duration) kv.DistributedLock {
	return &locker{
		key:   resource,
		ttl:   d,
		state: &sync.Mutex{},
	}
}
