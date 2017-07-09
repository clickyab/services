package aredis

import (
	"github.com/clickyab/services/assert"

	"context"

	"github.com/clickyab/services/initializer"

	"strings"

	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/clickyab/services/healthz"
	redis "gopkg.in/redis.v5"
)

var (
	// Client the actual pool to use with redis
	Client *redis.Client
)

type initRedis struct {
}

// Healthy return true if the databases are ok and ready for ping
func (initRedis) Healthy(context.Context) error {
	ping, err := Client.Ping().Result()
	if err != nil || strings.ToUpper(ping) == "PONG" {
		return fmt.Errorf("Redis PING failed. result was '%s' and the error was %s", ping, err)
	}

	return nil
}

// Initialize try to create a redis pool
func (i *initRedis) Initialize(ctx context.Context) {
	Client = redis.NewClient(
		&redis.Options{
			Network:  network.String(),
			Addr:     address.String(),
			Password: password.String(),
			PoolSize: poolsize.Int(),
			DB:       db.Int(),
		},
	)
	// PING the server to make sure every thing is fine
	assert.Nil(Client.Ping().Err())
	healthz.Register(i)
	logrus.Debug("redis is ready.")
	go func() {
		c := ctx.Done()
		assert.NotNil(c, "[BUG] context has no mean to cancel/deadline/timeout")
		<-c
		assert.Nil(Client.Close())
		logrus.Debug("redis finalized.")
	}()
}

func init() {
	initializer.Register(&initRedis{}, 0)
}
