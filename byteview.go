package geche

// A ByteView holds an immutable view of bytes.
// 抽象出来的缓存值类型
type ByteView struct {
	b []byte // 为了支持任意数据类型的存储，字符串、图片等
}

// Len
// 实现 Value 接口
func (bv ByteView) Len() int {
	return len(bv.b)
}

// String returns the data as a string, making a copy if necessary.
// 返回存储值的 string 形式
func (bv ByteView) String() string {
	return string(bv.b)
}

// ByteSlice returns a copy of the data as a byte slice.
// 返回 copy 的缓存值，防止缓存值被外部程序修改
func (bv ByteView) ByteSlice() []byte {
	return cloneBytes(bv.b)
}

// 拷贝缓存值
func cloneBytes(b []byte) []byte {
	c := make([]byte, 0)
	copy(c, b)
	return c
}
