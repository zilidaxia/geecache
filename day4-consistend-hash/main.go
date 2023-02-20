package main

import (
	"fmt"
	"log"
	geecache "lru.go/day3-http-server/geecache/lru"
	"net/http"
)

var db = map[string]string {
	"Tom": "630",
	"Jack": "589",
	"Sam": "567",
}
/*func main() {
	//回调函数 构建缓存
	geecache.NewGroup("scores",2<<10,geecache.GetterFunc(
		func(key string) ([]byte,error) {
			log.Println("[SlowDB] search key",key)
			if v,ok := db[key];ok {
				return []byte(v),nil
			}
			//缓存里面没有数据
			return nil,fmt.Errorf("%s not exist",key)
		}))
	addr := "localhost:9900"
	peers :=geecache.NewHTTPPool(addr)
	log.Println("geecache is running at",addr)
	log.Fatal(http.ListenAndServe(addr,peers))
}*/
func main() {
	geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := geecache.NewHTTPPool(addr)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}

