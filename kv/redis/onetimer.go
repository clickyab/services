package redis

import (
	"time"

	"github.com/clickyab/services/aredis"
)

type oneTimer struct {
	d   time.Duration
	key string
}

func (ot *oneTimer) Key() string {
	return ot.key
}

func (ot *oneTimer) Set(s string) string {
	b := aredis.Client.SetNX(ot.key, s, ot.d)
	if b.Err() != nil {
		return s
	}
	if b.Val() { // means it set the value so its for the first time
		return s
	}

	v := aredis.Client.Get(ot.key)
	if v.Err() != nil {
		return s
	}
	_ = aredis.Client.Expire(ot.key, ot.d) // ignore error
	return v.Val()
}
