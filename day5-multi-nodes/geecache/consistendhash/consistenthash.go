package consistendhash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

//定义函数类型Hash 用于计算hash值 这里使用自带的
type Hash func(data []byte) uint32

//Map 四个成员变量
type Map struct {
	hash Hash //hash函数
	replicas int //虚拟节点的倍数
	keys []int //哈希环
	hashMap map[int]string //虚拟节点和真实节点的映射
}
//构造函数  允许自定义虚拟节点倍数和Hash函数
func New(replicas int,fn Hash) *Map {
	m:= &Map{
		replicas:replicas,
		hash: fn,
		hashMap:make(map[int]string),
	}
	if m.hash == nil {
		m.hash=crc32.ChecksumIEEE
	}
	return m
}

//添加节点  任务 1.计算hash值节2.虚拟节点加入环 map映射真实节点和虚拟节点
//允许一次添加多个节点
func (m *Map) Add(keys ...string) {
	for _,key := range keys {
		//加入对应虚拟节点  虚拟节点的名称是i+key
		for i:=0;i<m.replicas;i++ {
			hash:=int(m.hash([]byte(strconv.Itoa(i)+key)))
			//hash值加入hash环中
			m.keys=append(m.keys,hash)
			//建立映射关系
			m.hashMap[hash] = key
		}
		//排一下序 好取值
		sort.Ints(m.keys)
	}
}

//获取Get
//步骤1.根据key计算hash值 2.顺时针找到第一个匹配的虚拟节点下标  3.通过hashmao映射得到真实的节点
func (m *Map) Get(key string) string{
	if len(m.keys) ==0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	//顺时针找第一个下标
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i]>=hash
	})
	//如果大于当前hash值的数不存在 即代表对应的索引是环开始的第一个节点 m.keys[0]
	//需要注意去一下余数 如果idx == len(m.keys)
	return m.hashMap[m.keys[idx%len(m.keys)]]
}

