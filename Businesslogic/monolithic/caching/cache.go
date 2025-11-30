package caching

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"karix.com/monolith/schemas"
)

type Cacheservice struct {
	rcli    *redis.Client
	local   *LocalCache
	ttl     time.Duration
	redisOk bool
}

func NewCacheService(rcli *redis.Client, local *LocalCache, ttl time.Duration) *Cacheservice {

	cs := &Cacheservice{
		rcli:    rcli,
		local:   local,
		ttl:     ttl,
		redisOk: false,
	}

	_, err := rcli.Ping().Result()

	if err == nil {
		cs.redisOk = true
	}

	return cs
}

func (cs *Cacheservice) GetUser(ctx context.Context, key string) (*schemas.User, bool, error) {

	if v, ok := cs.local.Get(key); ok {
		if user, ok := v.(*schemas.User); ok {
			return user, true, nil
		}
	}

	if cs.redisOk {
		val, err := cs.rcli.Get(key).Result()
		if err == nil {
			var user schemas.User
			if err := json.Unmarshal([]byte(val), &user); err != nil {
				cs.local.Set(key, nil, time.Minute*5)
				return nil, false, err
			}

			return &user, true, nil
		} else if err != nil {
			cs.redisOk = false
			return nil, false, err
		}

	}
	return nil, false, nil

}

func (cs *Cacheservice) SetUser(ctx context.Context, key string, user *schemas.User) error {

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	if cs.redisOk {
		err := cs.rcli.Set(key, data, cs.ttl).Err()
		if err != nil {
			cs.redisOk = false
			return err
		}
	}

	cs.local.Set(key, user, cs.ttl)
	return nil
}

type localItem struct {
	value  interface{}
	expiry time.Time
}

type LocalCache struct {
	m  map[string]*localItem
	mu sync.RWMutex
}

func NewLocalCache(defualtTTL time.Duration) *LocalCache {

	lc := &LocalCache{
		m: make(map[string]*localItem),
	}
	go lc.reaper()
	return lc
}

func (lc *LocalCache) Set(key string, value interface{}, ttl time.Duration) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.m[key] = &localItem{
		value:  value,
		expiry: time.Now().Add(ttl),
	}
}

func (lc *LocalCache) Get(key string) (interface{}, bool) {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	item, ok := lc.m[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(item.expiry) {
		return nil, false
	}
	return item.value, true
}

func (lc *LocalCache) reaper() {

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		lc.mu.Lock()
		for k, item := range lc.m {
			if item.expiry.Before(now) {
				delete(lc.m, k)
			}
		}
		lc.mu.Unlock()
	}

}
