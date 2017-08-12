package redlock

import (
	"services/dlock"
	"services/random"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"

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
}

const tryCooldown time.Duration = time.Millisecond * 200

// Lock is used to set a record in redis and tries until it gets its goal
func (m *mux) Lock() {
	m.locker.Lock()
	defer m.locker.Unlock()

	m.now = time.Now()

	m.value = <-random.ID

	for i := 0; i < m.retries; i++ {
		res := aredis.Client.SetNX(m.resource, m.value, m.TTL())
		if ok, err := res.Result(); ok == false || err != nil {
			time.Sleep(tryCooldown)
			continue
		}
		break
	}
}

// Unlock tries to get the record from redis and tries until it can
func (m *mux) Unlock() {
	m.locker.Lock()
	defer m.locker.Unlock()

	if m.value == "" {
		panic("unlocked before locking")
	}

	h := unlockScript
	for i := 0; i < m.retries; i++ {
		if time.Now().Sub(m.now) > m.TTL() {
			return
		}

		cmd := aredis.Client.Eval(h, []string{m.value})
		if cmd.Err() == nil {
			if val, ok := cmd.Val().(string); val != m.value || ok == false {
				time.Sleep(tryCooldown)
				continue
			}
		} else {
			logrus.Debug(cmd.Err().Error())
			time.Sleep(tryCooldown)
			continue
		}
		m.value = ""
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

// NewDistributedLock returns interface of a redlock
func NewDistributedLock(resource string, ttl time.Duration) dlock.DistributedLock {
	return &mux{
		retries:  int((ttl / tryCooldown).Nanoseconds() / 2),
		resource: resource,
		ttl:      ttl,
	}
}

func init() {
	dlock.Register(NewDistributedLock)
}
