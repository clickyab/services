package mock

import (
	"time"

	"github.com/clickyab/services/kv"
)

// AKiwi is a mock kiwi
type AKiwi struct {
	MasterKey string
	Data      map[string]int64
	Duration  time.Duration
}

var (
	atomicPre   = make(map[string]map[string]int64)
	atomicStore = make(map[string]*AKiwi)
)

// Drop the key
func (m *AKiwi) Drop() error {
	m.Data = make(map[string]int64)
	pre[m.MasterKey] = nil

	return nil
}

// Key return the parent key
func (m *AKiwi) Key() string {
	return m.MasterKey
}

// IncSubKey for increasing sub key
func (m *AKiwi) IncSubKey(key string, value int64) kv.AKiwi {
	t := m.Data[key]
	m.Data[key] = t + value
	return m
}

// DecSubKey for decreasing sub key
func (m *AKiwi) DecSubKey(key string, value int64) kv.AKiwi {
	return m.IncSubKey(key, value*-1)
}

// SubKey return a key
func (m *AKiwi) SubKey(key string) int64 {
	return m.Data[key]
}

// AllKeys from the store
func (m *AKiwi) AllKeys() map[string]int64 {
	return m.Data
}

// Save the entire keys (mostly first time)
func (m *AKiwi) Save(t time.Duration) error {
	m.Duration = t
	return nil
}

// TTL return the expiration time of this
func (m *AKiwi) TTL() time.Duration {
	return m.Duration
}

// NewAtomicMockStore is the new mock store
func NewAtomicMockStore(key string) kv.AKiwi {
	if k, ok := atomicStore[key]; ok {
		return k
	}
	var (
		data map[string]int64
		ok   bool
	)
	if data, ok = atomicPre[key]; !ok {
		data = make(map[string]int64)
	}
	m := &AKiwi{
		MasterKey: key,
		Data:      data,
	}

	atomicStore[key] = m
	return m
}

// SetAtomicMockData try to set mock data if needed
func SetAtomicMockData(key string, data map[string]int64) {
	atomicPre[key] = data
}

// GetAtomicMockStore is a function to get the mock store for testing
func GetAtomicMockStore() map[string]map[string]int64 {
	res := make(map[string]map[string]int64)
	for i := range store {
		res[i] = atomicStore[i].Data
	}

	return res
}

// ResetEaev the entire mock
func ResetEaev() {
	atomicStore = make(map[string]*AKiwi)
}
