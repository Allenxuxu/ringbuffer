package ringbuffer

import (
	"strings"
	"testing"
)

func BenchmarkRingBuffer_Sync(b *testing.B) {
	rb := New(1024)
	data := []byte(strings.Repeat("a", 512))
	buf := make([]byte, 512)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Write(data)
		rb.Read(buf)
	}
}
