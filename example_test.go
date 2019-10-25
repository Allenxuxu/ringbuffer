package ringbuffer

import "fmt"

func ExampleRingBuffer() {
	rb := New(1024)
	_, _ = rb.Write([]byte("abcd"))
	fmt.Println(rb.Length())
	fmt.Println(rb.free())
	buf := make([]byte, 4)

	_, _ = rb.Read(buf)
	fmt.Println(string(buf))
	// Output: 4
	// 1020
	// abcd
}
