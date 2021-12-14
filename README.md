# ringbuffer

自动扩容的循环缓冲区实现

[![Github Actions](https://github.com/Allenxuxu/ringbuffer/workflows/CI/badge.svg)](https://github.com/Allenxuxu/ringbuffer/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/Allenxuxu/ringbuffer)](https://goreportcard.com/report/github.com/Allenxuxu/ringbuffer)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/eb54ecf0096244d39949843efb895918)](https://www.codacy.com/manual/Allenxuxu/ringbuffer?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=Allenxuxu/ringbuffer&amp;utm_campaign=Badge_Grade)
[![GoDoc](https://godoc.org/github.com/Allenxuxu/ringbuffer?status.svg)](https://godoc.org/github.com/Allenxuxu/ringbuffer)
[![LICENSE](https://img.shields.io/badge/LICENSE-MIT-blue)](https://github.com/Allenxuxu/ringbuffer/blob/master/LICENSE)
[![Code Size](https://img.shields.io/github/languages/code-size/Allenxuxu/ringbuffer.svg?style=flat)](https://img.shields.io/github/languages/code-size/Allenxuxu/ringbuffer.svg?style=flat)

## 使用

```go
package main

import (
	"fmt"
	"github.com/Allenxuxu/ringbuffer"
)

func main() {
	rb := ringbuffer.New(2)

	// 自动扩容
	fmt.Println(rb.Capacity())  //2
	fmt.Println(rb.Length())    //0

	rb.Write([]byte("ab"))
	fmt.Println(rb.Capacity())  //2
	fmt.Println(rb.Length())    //2

	rb.Write([]byte("cd"))
	fmt.Println(rb.Capacity())  //4
	fmt.Println(rb.Length())    //4

	// VirtualXXX 函数便捷操作
	rb = New(1024)
	_, _ = rb.Write([]byte("abcd"))
	fmt.Println(rb.Length())
	fmt.Println(rb.free())
	buf := make([]byte, 4)

	_, _ = rb.Read(buf)
	fmt.Println(string(buf))

	rb.Write([]byte("1234567890"))
	rb.VirtualRead(buf)
	fmt.Println(string(buf))
	fmt.Println(rb.Length())
	fmt.Println(rb.VirtualLength())
	rb.VirtualFlush()
	fmt.Println(rb.Length())
	fmt.Println(rb.VirtualLength())

	rb.VirtualRead(buf)
	fmt.Println(string(buf))
	fmt.Println(rb.Length())
	fmt.Println(rb.VirtualLength())
	rb.VirtualRevert()
	fmt.Println(rb.Length())
	fmt.Println(rb.VirtualLength())
	// Output: 4
	// 1020
	// abcd
	// 1234
	// 10
	// 6
	// 6
	// 6
	// 5678
	// 6
	// 2
	// 6
	// 6
}
```

## 扩容

ringbuffer底层存储结构为golang切片。当ringbuffer需要扩容时，扩容策略参考golang切片append策略:  
https://github.com/golang/go/blob/ac0ba6707c1655ea4316b41d06571a0303cc60eb/src/runtime/slice.go#L125  
1. 如果期望容量大于当前容量的两倍就会使用期望容量；
2. 如果当前切片的长度小于 1024 就会将容量翻倍；
3. 如果当前切片的长度大于 1024 就会每次增加 25% 的容量，直到新容量大于期望容量；

当调用ringbuffer.Reset()时，buffer将会缩容回初始化时的大小。

```go
rb := New(2)
fmt.Println(rb.Length())   // 0
fmt.Println(rb.Capacity()) // 2

rb.Write([]byte("abc"))
fmt.Println(rb.Length())   // 3
fmt.Println(rb.Capacity()) // 4

rb.Write([]byte(strings.Repeat("a", 1024)))
fmt.Println(rb.Length())   // 1027
fmt.Println(rb.Capacity()) // 1027

rb.WriteByte('a')
fmt.Println(rb.Length())   // 1028
fmt.Println(rb.Capacity()) // 1283

rb.Reset()
fmt.Println(rb.Length())   // 0
fmt.Println(rb.Capacity()) // 2
``` 


## 参考

https://github.com/smallnest/ringbuffer

## 感谢

- [李舒畅](https://github.com/MrChang0)
- [Harold2017](https://github.com/Harold2017)
