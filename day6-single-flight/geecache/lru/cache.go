package geecache

import (
	"lru.go/day1-lru/geecache/lru"
	"sync"
)

//定义cache结构体
type cache struct {
	//mu 互斥锁 lru
	mu  sync.Mutex
	//lru包下的cache结构  使用指针 来进行调用之后
	lru *lru.Cache
	//缓存大小 用于初始化缓存
	cacheBytes int64
}
//在add 和 get 方法内设置锁即可
func (c *cache) add(key string,value ByteView){
	c.mu.Lock()
	defer c.mu.Unlock()
	//首先判断lru是否初始化 没有给其初始化
	if c.lru == nil {
		c.lru=lru.New(c.cacheBytes,nil)
	}
	//value就是缓存数据  这里传入之后胡 lru中的add添加到map里面 key value后面在value其实就是当前k缓存大小
	c.lru.Add(key,value)
}
func (c *cache) get(key string) (value ByteView,ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	//lru中查找key
	if v,ok := c.lru.Get(key);ok {
		return v.(ByteView), ok
	}
	return
}
