package redis

import (
	"sync"
	"time"

	"sync/atomic"

	"github.com/clickyab/services/aredis"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/random"
)

const unlockScript = `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        else
            return 0
        end
        `

type mux struct {
	locker sync.Mutex

	ttl      time.Duration
	now      time.Time
	resource string

	value   string
	retries int

	swap int32
}

var tryCoolDown = config.RegisterDuration("services.kv.redlock.cooldown", time.Millisecond*200, "cooldown for redlock algorithm in redis")

// TODO : this is not compatible with redis cluster. it work only when there is one redis instance

// Lock is used to set a record in redis and tries until it gets its goal
func (m *mux) Lock() {
	// Make sure this is locked here, and concurrent call to this
	// is blocked by the simple mutex.
	// this is a distributed lock, and also it means it must block the same
	// call in the same routine too!
	m.locker.Lock()

	assert.True(atomic.CompareAndSwapInt32(&m.swap, 0, 1), "invalid value on lock")
	m.now = time.Now()

	m.value = <-random.ID

	for i := 0; i < m.retries; i++ {
		res := aredis.Client.SetNX(m.resource, m.value, m.TTL())
		if ok, err := res.Result(); ok == false || err != nil {
			time.Sleep(tryCoolDown.Duration())
			continue
		}
		break
	}
}

// Unlock tries to get the record from redis and tries until it can
func (m *mux) Unlock() {
	// must be unlocked after the call, this block concurrent call in same process
	defer m.locker.Unlock()

	assert.True(atomic.LoadInt32(&m.swap) == 1, "not locked?")
	defer func() {
		assert.True(atomic.CompareAndSwapInt32(&m.swap, 1, 0), "not locked?")
	}()

	h := unlockScript
	for i := 0; i < m.retries; i++ {
		if time.Since(m.now) > m.TTL() {
			return
		}

		cmd := aredis.Client.Eval(h, []string{m.value})
		if cmd.Err() == nil {
			if val, ok := cmd.Val().(string); val != m.value || ok == false {
				time.Sleep(tryCoolDown.Duration())
				continue
			}
		} else {
			time.Sleep(tryCoolDown.Duration())
			continue
		}
		break
	}
}

// Resource returns resource for no reason
func (m *mux) Resource() string {
	return m.resource
}

// TTL returns the duration of a lock
func (m *mux) TTL() time.Duration {
	return m.ttl
}

// newRedisDistributedLock returns interface of a redlock
func newRedisDistributedLock(resource string, ttl time.Duration) kv.DistributedLock {
	return &mux{
		retries:  int((ttl / tryCoolDown.Duration()).Nanoseconds() / 2),
		resource: resource,
		ttl:      ttl,
	}
}
