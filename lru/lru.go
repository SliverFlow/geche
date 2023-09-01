package lru

import (
	"container/list"
	"log"
)

// 淘汰算法选择 （LFU 最近最少使用算法）
// 实现：
// 维护链表存储数据的最大空间以及当前使用空间
// 使用标准库中的双向链表存储， 节点上存储数据值 （键值对）
// 维护数据值的键 和 链表上节点指针的 map
// 提供可供用户操作的回调函数
// 约定（list.list front 为队尾 list.list back 为队头）
// 队尾：最近最常使用 (front)
// 队头：最不常使用 (back)

type (
	// OnEvictedFunc 回调函数类型
	OnEvictedFunc func(key string, value Value)

	// Cache 核心存储结构
	Cache struct {
		maxBytes  int64                    // 最大存储空间 实际字节为单位
		nBytes    int64                    // 现使用空间大小
		ll        *list.List               // 标准库实现的双向链表
		hashMap   map[string]*list.Element // 键为 string ,值为 双向链表上对应节点的指针
		OnEvicted OnEvictedFunc            // 缓存值被删除时的回调 可为 nil
	}

	// Value 存储键值对的值类型需实现这个接口
	// 后续方便计算现使用空间大小
	Value interface {
		Len() int
	}

	// 存储值的数据结构 键值对
	entry struct {
		key   string
		value Value
	}
)

// New is the Constructor of Cache
func New(maxBytes int64, evictedFunc OnEvictedFunc) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		hashMap:   make(map[string]*list.Element),
		OnEvicted: evictedFunc,
	}
}

// Get look ups a key's value
// 将元素移动至队尾 约定 头部为队尾
func (c *Cache) Get(key string) (value Value, ok bool) {
	if key == "" {
		return nil, false
	}
	if ele, ok := c.hashMap[key]; ok { // hashMap 中存储了改键的节点指针
		c.ll.MoveToFront(ele)    // 将此节点移动至链表头 由于是双向链表 队尾和队头是相对的，这里约定 Front 为队尾
		kv := ele.Value.(*entry) // 由于存储值时接收类型为 any ,所以取值是需要断言为自定义的存储值数据类型
		return kv.value, ok
	}
	return nil, false
}

// Add adds or update a value to the cache.
// hashMap 中无对应的键为添加 反之 则为更新
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.hashMap[key]; ok { // 更新
		c.ll.MoveToFront(ele) // 元素移动至队列头
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len()) // 重新计算现使用空间大小
		log.Printf("update itme key [%s] old value [%v] value [%v]", key, kv.value, value)
		kv.value = value // 更新值
	} else { // 添加
		ele := c.ll.PushFront(&entry{
			key:   key,
			value: value,
		}) // 存储节点
		c.hashMap[key] = ele                             // 更新 key 和 节点指针映射关系
		c.nBytes += int64(len(key)) + int64(value.Len()) // 重新计算现使用空间大小
		log.Printf("add itme key [%s] value [%v]", key, value)
	}
	// 循环删除节点 一直到现使用空间大小小于最大存储空间大小
	for c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.removeOldest()
	}
}

// RemoveOldest removes the oldest item
// 不通过用户控制
// 实际是淘汰算法的体现
func (c *Cache) removeOldest() {
	ele := c.ll.Back() // 队列的最后一个元素
	if ele != nil {
		c.ll.Remove(ele) // 移除元素
		kv := ele.Value.(*entry)
		delete(c.hashMap, kv.key)                              // 删除维护 key 和 节点指针映射关系
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len()) // 重新计算现使用空间大小
		if c.OnEvicted != nil {                                // 执行回调
			c.OnEvicted(kv.key, kv.value)
		}
		log.Printf("remove oldest itme key [%s]", kv.key)
	}
}

// Len the number of cache entries
// 方便测试 返回当前存储的节点个数
func (c *Cache) Len() int {
	return c.ll.Len()
}

// Keys the all key of cache list
func (c *Cache) Keys() *[]string {
	keys := make([]string, 0)
	for k, _ := range c.hashMap {
		keys = append(keys, k)
	}
	return &keys
}
