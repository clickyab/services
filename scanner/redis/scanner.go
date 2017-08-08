package redis

import (
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/redis"
	"github.com/clickyab/services/scanner"
)

type counter struct {
	pattern string
	count   uint64
}

// Next iterates on more keys
func (c counter) Next() ([]string, bool) {
	scanCmd := aredis.Client.Scan(c.count, c.pattern, 10)
	assert.Nil(scanCmd.Err())

	var keys []string
	keys, c.count = scanCmd.Val()
	if c.count == 0 {
		return keys, true
	}

	return keys, false
}

func newScanner(pattern string) scanner.Scanner {
	return counter{pattern: pattern}
}

func init() {
	scanner.Register(newScanner)
}
