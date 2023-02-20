package geecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

//创建结构体HTTPPool，作为服务端和客户端通信的节点

//basepath
const defaultBasePath = "/_geecache/"
type HTTPPool struct {
	self      string  //记录自己的地址
	basePath  string
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








