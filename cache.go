package geche

import (
	"geche/lru"
	"sync"
)

// 为 lru 添加并发支持
// 避免并发情况下导致资源冲突问题
type cache struct {
	mu         sync.Mutex // 锁
	lru        *lru.Cache // 淘汰算法
	cacheBytes int64      // 缓存空间大小
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil) // 延迟实例化 lru
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), true
	}
	return
}
