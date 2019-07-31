# ringbuffer
自动扩容的循环缓冲区实现

- RingBuffer 非线程安全
- AtomicRingBuffer 线程安全

### 使用

```go
package main

import (
	"fmt"
	"github.com/Allenxuxu/ringbuffer"
)

func main() {
	rb := ringbuffer.New(2)
  //rb := ringbuffer.NewAtomicRingBuffer(2)
	fmt.Println(rb.Capacity())  //2
	fmt.Println(rb.Length())    //0

	rb.Write([]byte("ab"))
	fmt.Println(rb.Capacity())  //2
	fmt.Println(rb.Length())    //2

	rb.Write([]byte("cd"))
	fmt.Println(rb.Capacity())  //4
	fmt.Println(rb.Length())    //4
}
```

### 参考
https://github.com/smallnest/ringbuffer
