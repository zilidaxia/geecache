package geecache

import (
	"fmt"
	"log"
	"sync"
)

//核心数据结构Group
type Group struct {
	name string
	//回调函数
	getter Getter
	mainCache cache
}

var (
	mu sync.Mutex
	//map 存group  通过key value形式 比如学生信息  学生成绩
	groups =make(map[string]*Group)
)
//NewGroup
func NewGroup(name string,cacheBytes int64,getter Getter) *Group{
	//构建时需要加锁 类似redis中的排它锁构建缓存
	if getter==nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes:cacheBytes},
	}
	//map更新一下
	groups[name] = g
	return g
}
//GetGroup
func GetGroup(name string) *Group {
	mu.Lock()
	defer mu.Unlock()
	g:=groups[name]
	return g
}

//首先Getter 是一个接口 其含有Get函数  当我们缓存未命中 需要从数据库拿到数据时
//Getter
type Getter interface {
	Get(key string) ([]byte,error)
}
//定义一个函数类型 其类型跟上面的Get一样
type GetterFunc func(key string) ([]byte,error)

//回调函数
//这里 函数类型实现某一个接口，我们传入是方便传入不同类型的数据库作为参数 类似多态的盖帘
func (f GetterFunc) Get(key string) ([]byte,error) {
	return f(key)
}

//Get方法具体实现
//1.去请求一个缓存 两种情况 缓存存在获取 缓存不存在去数据库拿数据
func (g *Group) Get(key string) (ByteView, error) {
	if key=="" {
		return ByteView{},fmt.Errorf("key is required")
	}
	//获取缓存
	if v,ok := g.mainCache.get(key); ok{
		log.Println("[GeeCache] hit")
		return v,nil
	}
	//构建缓存
	return g.load(key)
}
func (g *Group) load(key string) (value ByteView,err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView,error) {
	bytes,err:=g.getter.Get(key)
	if err != nil {
		return ByteView{},err
	}
	//数据拷贝下来 变为只读模式
	value := ByteView{b: cloneBytes(bytes)}
	//添加到缓存里面
	g.populateCache(key,value)
	return value, nil
}

func (g *Group) populateCache(key string,value ByteView) {
	g.mainCache.add(key,value)
}