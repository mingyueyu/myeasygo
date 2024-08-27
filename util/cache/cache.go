package cache

import (
	"fmt"
	"sync"
	"time"
)

var c = NewCache()

// CacheItem 代表缓存中的一项
type CacheItem struct {
	Value      interface{}
	Expiration int64
}

// Cache 全局缓存
type Cache struct {
	items sync.Map
}

// NewCache 创建一个新的 Cache 实例
func NewCache() *Cache {
	return &Cache{}
}

// Set 将数据存入缓存
func Set(key string, value interface{}) {
	expiration := time.Now().Add(2 * time.Hour).UnixNano()
	c.items.Store(key, CacheItem{
		Value:      value,
		Expiration: expiration,
	})
}

// 设置缓存，缓存有效期几秒
func SetDuration(key string, value interface{}, second int64) {
	duration := time.Duration(second) * time.Second
	expiration := time.Now().Add(duration).UnixNano()
	c.items.Store(key, CacheItem{
		Value:      value,
		Expiration: expiration,
	})
}

// Get 从缓存中获取数据
func Get(key string) (interface{}, bool) {
	item, found := c.items.Load(key)
	if !found {
		return nil, false
	}

	cacheItem := item.(CacheItem)
	if time.Now().UnixNano() > cacheItem.Expiration {
		c.items.Delete(key)
		return nil, false
	}

	return cacheItem.Value, true
}

// Delete 从缓存中删除数据
func Delete(key string) {
	c.items.Delete(key)
}

// Cleanup 清理过期的缓存
func Cleanup() {
	now := time.Now().UnixNano()
	c.items.Range(func(key, value interface{}) bool {
		cacheItem := value.(CacheItem)
		if now > cacheItem.Expiration {
			fmt.Println(key, ":", value)
			c.items.Delete(key)
		}
		return true
	})
}
