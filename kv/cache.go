package kv

import (
	"encoding/gob"
	"io"
	"time"
)

// Serializable represent the object that can be serialized
type Serializable interface {
	// Decode is the decoder of this function
	Decode(io.Writer) error
	// Encode is the encoder function
	Encode(io.Reader) error
}

// CacheProvider is the cacheFactory backend
type CacheProvider interface {
	// Do is called to store the cacheFactory
	Do(string, Serializable, time.Duration) error
	// Hit called when we need to load the cacheFactory
	Hit(string, Serializable) error
}

// CacheWrapper is a provider with support for inner entity
type CacheWrapper interface {
	Serializable
	// Entity return the cached object
	Entity() interface{}
}

type cachable struct {
	entity interface{}
}

var cacheFactory CacheProvider

// Decode try to decode cookie profile into gob
func (cp *cachable) Decode(w io.Writer) error {
	enc := gob.NewEncoder(w)
	return enc.Encode(cp.entity)
}

// Encode try to encode cookie profile from gob
func (cp *cachable) Encode(i io.Reader) error {
	dnc := gob.NewDecoder(i)
	return dnc.Decode(cp.entity)
}

func (cp *cachable) Entity() interface{} {
	return cp.entity
}

// Do the entity
func Do(k string, e Serializable, t time.Duration, err error) error {
	if err != nil {
		return err
	}
	regLock.RLock()
	defer regLock.RUnlock()

	return cacheFactory.Do(k, e, t)
}

// Hit the cacheFactory
func Hit(key string, out Serializable) error {
	regLock.RLock()
	defer regLock.RUnlock()

	return cacheFactory.Hit(key, out)
}

// CreateWrapper return an cachable object for this ntt
func CreateWrapper(ntt interface{}) CacheWrapper {
	return &cachable{
		entity: ntt,
	}
}
