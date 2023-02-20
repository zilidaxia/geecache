package geecache
//b是一个只读的数据结构，存储真实的缓存值 byte是为了存任意资源 比如图片等
type ByteView struct {
	b []byte
}
//实现一个len方法  返回缓存所占内存大小
func (v ByteView) Len() int{
	return len(v.b)
}
//b是只读的 使用byteslice方法返回一个拷贝 防止缓存在被外部程序修改
func (v ByteView) ByteSlice() []byte{
	return cloneBytes(v.b)
}
func (v ByteView) String() string {
	return string(v.b)
}
func cloneBytes(b []byte) []byte {
	c := make([]byte,len(b))
	copy(c,b)
	return c
}