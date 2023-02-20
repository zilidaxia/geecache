package lru

import (
	"testing"
)

type String string
func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T){
	lru := New(int64(0),nil)
	//fmt.Println(lru.nbyBytes+1)
	lru.Add("key1",String("1234"))
	//判断一下是否找到对应key
	if v,ok:=lru.Get("key1");!ok || string(v.(String))!="1234"{
		t.Fatal("cache hit key1=1234 failed")
	}
	if _,ok:=lru.Get("key2");ok{
		t.Fatalf("cache miss key2 failed")
	}
}
//当内存超过设定值时，测试是否会触发无用节点
func TestRemoveOldest(t *testing.T){
	k1,k2,k3 := "key1","key2","k3"
	v1,v2,v3 := "value1","value2","v3"
	cap := len(k1+k2+v1+v2)
	lru:=New(int64(cap),nil)
	lru.Add(k1,String(v1))
	lru.Add(k2,String(v2))
	lru.Add(k3,String(v3))
	if _, ok := lru.Get("key1");ok||lru.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}
