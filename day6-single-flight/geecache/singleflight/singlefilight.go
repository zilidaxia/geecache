package singleflight

import "sync"

//call代表一次HTTP 客户端请求
type call struct {
	wg sync.WaitGroup
	val interface{}
	err error
}
//singleflight 的主要结构group  用来管理不同key的请求
type Group struct {
	//为保证m 初始化 不被并发读写
	mu sync.Mutex
	m map[string]*call
}

//核心Do方法
//注意此时这里的Group是为啦管理不同key请求 而不是之前的Group
func (g *Group) Do(key string,fn func() (interface{},error)) (interface{},error) {
	//延时初始化m
	g.mu.Lock()
	if g.m==nil {
		g.m=make(map[string]*call)
	}
	//获取在m 中key对应请求  如果获取成功 代表之前已经有个当前key请求
	if c,ok :=g.m[key]; ok {
		//首先把mu解锁
		g.mu.Unlock()
		//进行同步操作  需要等第一个请求（线程）完成
		c.wg.Wait()
		//待第一个请求完成 获取其返回值
		return c.val,c.err
	}
	//接下来这部分处理如果是第一个请求
	c :=new(call)
	//Add添加一个 线程开启提示
	c.wg.Add(1)
	//key放入map中
	g.m[key]=c
	//也就是对map key读写操作 需要互斥锁实现  对请求限制是需要同步wait实现
	g.mu.Unlock()
	//调用请求函数获取返回值 使用call对象接受
	c.val,c.err=fn()
	//做完这个线程任务需要通知其他请求（线程）
	c.wg.Done()
	//接下来又进行对map删除key操作  因为此时请求已经完成 必须删除key
	//防止比如下一个同一个key来 一位已经有一次key请求  无限等待
	//互斥锁加上
	g.mu.Lock()
	delete(g.m,key)
	g.mu.Unlock()
	//最后返回请求函数的结果即可
	return c.val,c.err
}

