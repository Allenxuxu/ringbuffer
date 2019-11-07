# ringbuffer

自动扩容的循环缓冲区实现

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

## 参考

https://github.com/smallnest/ringbuffer

## 感谢

- [李舒畅](https://github.com/MrChang0)
- [Harold2017](https://github.com/Harold2017)
