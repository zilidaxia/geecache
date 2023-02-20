package geecache

import (
	"fmt"
	"io/ioutil"
	"log"
	"lru.go/day5-multi-nodes/geecache/consistendhash"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

//创建结构体HTTPPool，作为服务端和客户端通信的节点

//basepath
const (
	defaultBasePath = "/_geecache/"
	defaultReplicas = 50
)
type HTTPPool struct {
	self      string  //记录自己的地址
	basePath  string
	mu 		  sync.Mutex
	//根据key来选择节点
	peers     *consistendhash.Map
	//根据url
	httpGetter map[string]*httpGetter
}

//New
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self: self,
		basePath:defaultBasePath,
	}
}
//日志输出服务端名字
func (p *HTTPPool) Log(format string,v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}
//实现ServerHTTP方法
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter,r *http.Request) {
	//1.快速判断url地址的前缀是否满足base
	if !strings.HasPrefix(r.URL.Path,p.basePath) {
		panic("HTTPPool serving unexpectd path "+ r.URL.Path)
	}
	p.Log("%s %s",r.Method,r.URL.Path)
	//约定访问格式是 /basePath/groupname/key
	//分割路径
	parts := strings.SplitN(r.URL.Path[len(p.basePath):],"/",2)
	//判断path中格式是否正确
	if len(parts) != 2 {
		http.Error(w,"bad request",http.StatusBadRequest)
		return
	}
	//获取对应group 之后获取对应name
	grouName := parts[0]
	key := parts[1]
	//通过GetGroup获取group
	group:=GetGroup(grouName)
	//group==nil 代表没找到
	if group == nil {
		http.Error(w,"no such group: "+grouName,http.StatusNotFound)
		return
	}
	//获取对应key
	view,err := group.Get(key)
	if err != nil {
		fmt.Println("没找到")
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	//将获取的数据同response写回即可
	//先封装头部
	w.Header().Set("Content-Type","application/octet-stream")
	//这里返回给其拷贝体 其他端就没办法修改
	w.Write(view.ByteSlice())
}

//创建HTTPGetter结构体
type httpGetter struct {
	baseURL string
}
//实现Get方法
func (h *httpGetter) Get(group string,key string) ([]byte,error) {
	//url路径封装一下格式 base/group/key
	 u := fmt.Sprintf(
	 	"%v%v/%v",
	 	h.baseURL,
	 	url.QueryEscape(group),
	 	url.QueryEscape(key),
	 	)
	 //发送请求  并接受消息
	 res, err :=http.Get(u)
	 if err!= nil {
	 	return nil,err
	 }
	 //最后需要关闭response
	 defer res.Body.Close()
	 //先对有异常情况进行处理 通过比对http响应的状态码
	 if res.StatusCode !=http.StatusOK {
	 	return nil,fmt.Errorf("sercer returned:%v",res.Status)
	 }
	 //拿数据 从res中获取
	 bytes,err := ioutil.ReadAll(res.Body)
	 if err != nil {
	 	return nil,fmt.Errorf("reading response body: %v",err)
	 }
	 return bytes,nil
}
var _PeerGetter = (*httpGetter)(nil)
//Set方法传入节点 实例化hash算法  为每一个节点创建一个HTTP客户端httpGetter
func (p *HTTPPool) Set(peers ...string) {
	//传入节点时 上锁
	p.mu.Lock()
	defer mu.Unlock()
	p.peers = consistendhash.New(defaultReplicas,nil)
	//将每个节点传入Add中添加节点
	p.peers.Add(peers...)
	//这里就是给每一个节点 分配一个HTTP客户端
	p.httpGetter=make(map[string]*httpGetter,len(peers))
	for _,peer :=range peers {
		//每一个客户端 分配一个ip地址
		p.httpGetter[peer] = &httpGetter{baseURL:peer+p.basePath}
	}
}

//通过传入key 选择节点 返回节点所对应的HTTP客户端
func (p *HTTPPool) PickPeer(key string) (PeerGetter,bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	//获取节点 通过get方法 获取失败get方法返回的是" " 而且要判定此时返回的不是当前节点
	//需要向远处节点获取数据
	if peer := p.peers.Get(key);peer != "" &&peer !=p.self{
		p.Log("Pick peer %s",peer)
		return p.httpGetter[peer],true
	};
	return nil,false
}






