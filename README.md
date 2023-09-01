### 目录结构：
```shell
example # 示例文件
  -- 
lru     
  -- lru.go # 淘汰算法
  -- lru_test.go # 淘汰算法测试文件
consistenthash
  -- consistenthash.go # 一致性 hash 
  -- consistenthash_test.go # 一致性 hash 测试文件
```

### Len 函数
其作用是用于计算数组(包括数组指针)、切片(slice)、map、channel、字符串等数据类型的长度，注意，结构休(struct)、整型布尔等不能作为参数传给len函数。
- 数组或数组指针：返回元素个数
- map和slice: 元素个数
- channel:通道中未读的元素个数
- 字符串：字节数，并非字符串的字符数
- 当V的值为nil值，len返回0