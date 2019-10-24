# ringbuffer
自动扩容的循环缓冲区实现

### 使用

```go
package main

import (
	"fmt"
	"github.com/Allenxuxu/ringbuffer"
)

func main() {
	rb := ringbuffer.New(2)

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

### 感谢

- [李舒畅](https://github.com/MrChang0)
- [Harold2017](https://github.com/Harold2017)
