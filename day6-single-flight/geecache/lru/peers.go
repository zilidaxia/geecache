package geecache

//两个接口
//通过传入的key  选择相应的节点的PeerGetter
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter,ok bool)
}

//通过传入的group找缓存值     这里PeerGetter对应HTTP客户端
type PeerGetter interface {
	Get(group string,key string) ([]byte,error)

}
