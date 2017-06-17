package redis

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/clickyab/services/cache"
	"github.com/clickyab/services/redis"
)

type redisCache struct {
}

// Sha1 is the sha1 generation func
func sha(k string) string {
	h := sha1.New()
	_, _ = h.Write([]byte(k))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Do is called to store the cache
func (redisCache) Do(e cache.Cacheable, t time.Duration) error {
	name := "CACHE_" + sha(e.String())
	target := &bytes.Buffer{}
	err := e.Decode(target)
	if err != nil {
		return err
	}

	res := aredis.Client.Set(name, target.String(), t)
	return res.Err()
}

// Hit called when we need to load the cache
func (redisCache) Hit(key string, e cache.Cacheable) error {
	name := "CACHE_" + sha(key)
	res := aredis.Client.Get(name)
	if err := res.Err(); err != nil {
		return err
	}
	data, err := res.Result()
	if err != nil {
		return err
	}
	buf := bytes.NewBufferString(data)
	return e.Encode(buf)
}

func init() {
	cache.Register(&redisCache{})
}
