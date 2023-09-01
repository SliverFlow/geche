package lru

import (
	"fmt"
	"strconv"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestLru(t *testing.T) {
	cache := New(int64(100), nil)
	for i := 0; i < 20; i++ {
		add(cache, "key"+strconv.Itoa(i), String("value"+strconv.Itoa(i)))
	}
	cache.Add("key19", String("我是自己添加的存储值"))
	fmt.Println("len", cache.Len())
	fmt.Println("keys", cache.Keys())
	value, _ := cache.Get("key17")
	fmt.Printf("GET KEY 【%s】value【%v】\n", "key17", value)
}

func add(c *Cache, key string, value Value) {
	c.Add(key, value)
}
