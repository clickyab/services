package redis

import (
	"sync"
	"time"

	"strconv"

	"github.com/clickyab/services/aredis"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/kv"
)

type atomicKiwiRedis struct {
	key  string
	v    map[string]int64
	lock sync.Mutex
}

func (kr *atomicKiwiRedis) TTL() time.Duration {
	d := aredis.Client.TTL(kr.key)
	r, _ := d.Result()

	return r
}

func (kr *atomicKiwiRedis) Drop() error {
	kr.lock.Lock()
	defer kr.lock.Unlock()

	kr.v = make(map[string]int64)
	d := aredis.Client.Del(kr.key)
	return d.Err()
}

// Key return the parent key
func (kr *atomicKiwiRedis) Key() string {
	return kr.key
}

// IncSubKey for increasing sub key
func (kr *atomicKiwiRedis) IncSubKey(key string, value int64) kv.AKiwi {
	res := aredis.Client.HIncrBy(kr.key, key, value)
	if res.Err() != nil {
		kr.v[key] = value
		return kr
	}
	r, err := res.Result()
	if err != nil {
		kr.v[key] = value
		return kr
	}
	kr.v[key] = r
	return kr
}

// IncSubKey for decreasing sub key
func (kr *atomicKiwiRedis) DecSubKey(key string, value int64) kv.AKiwi {
	return kr.IncSubKey(key, -value)
}

// SubKey return a key
func (kr *atomicKiwiRedis) SubKey(key string) int64 {
	kr.lock.Lock()
	defer kr.lock.Unlock()

	if v, ok := kr.v[key]; ok {
		return v
	}
	res := aredis.Client.HIncrBy(kr.key, key, 0)
	if res.Err() != nil {
		return 0
	}

	r, err := res.Result()
	if err != nil {
		return 0
	}

	return r
}

// AllKeys from the store
func (kr *atomicKiwiRedis) AllKeys() map[string]int64 {
	kr.v = map[string]int64{}
	res := aredis.Client.HGetAll(kr.key)

	if res.Err() != nil {
		return kr.v
	}

	r, err := res.Result()
	if err != nil {
		return kr.v
	}
	f := make(map[string]int64)
	for k, v := range r {
		tv, e := strconv.ParseInt(v, 10, 64)
		assert.Nil(e)
		f[k] = tv
	}
	kr.v = f
	return kr.v
}

// Save the entire keys (mostly first time)
func (kr *atomicKiwiRedis) Save(t time.Duration) error {
	b := aredis.Client.Expire(kr.key, t)
	return b.Err()
}

// NewRedisAEAVStore return a redis store for eav
func newRedisAEAVStore(key string) kv.AKiwi {
	return &atomicKiwiRedis{
		key:  key,
		v:    make(map[string]int64),
		lock: sync.Mutex{},
	}
}
