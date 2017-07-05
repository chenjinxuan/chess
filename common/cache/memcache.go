// go-cache is an in-memory key:value store/cache similar to memcached that is suitable for applications running on a single machine. Its major
// advantage is that, being essentially a thread-safe map[string]interface{} with expiration times, it doesn't need to serialize or transmit its
// contents over the network.
// github.com/patrickmn/go-cache
package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var Memcache *cache.Cache

func NewMemcache(defaultExpiration, cleanupInterval time.Duration) {
	if Memcache == nil {
		Memcache = cache.New(defaultExpiration, cleanupInterval)
	}
}
