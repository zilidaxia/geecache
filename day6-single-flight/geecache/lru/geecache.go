package geecache

import (
	"fmt"
	"log"
	"lru.go/day6-single-flight/geecache/singleflight"
	"sync"
)

//核心数据结构Group
type Group struct {
	name string
	//回调函数
	getter Getter
	mainCache cache
	//注入HTTPPool  peers  一群节点
	peers PeerPicker
	//加入singleflight 来实现防止缓存击穿
	loader *singleflight.Group
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
		loader:    &singleflight.Group{},
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

//将实现PeerPicker接口的HTTPPool注入到Group  实际上类似初始化
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers !=nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers=peers
}

//加载资源
//流程:1.首先看当前这个节点是否有缓存值 没有就去远程节点拿
func (g *Group) load(key string) (value ByteView,err error) {
	//将原来的代码逻辑放在do里面接口
	viewi, err :=g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok :=g.peers.PickPeer(key);ok {
				if value ,err =g.getFromPeer(peer,key);err==nil {
					return value,nil
				}
				//远程节点没访问到数据
				log.Println("GeeCache Failed to get from peer",err)
			}
		}
		//本地获取 构建缓存
		return g.getLocally(key)
	})
	if err != nil {
		return viewi.(ByteView),nil
	}
	return
}

//使用啦实现PeerGetter接口的httpGetter从
//从远程节点获取资源
func (g *Group) getFromPeer(peer PeerGetter,key string) (ByteView,error) {
	//由于之前httpgetter 实现啦PeerGetter接口的Get方法
	//通过远程节点的访问 HTTP通信  获取对应缓存
	bytes,err := peer.Get(g.name,key)
	if err != nil {
		return ByteView{},err
	}
	return ByteView{b:bytes},nil
}









