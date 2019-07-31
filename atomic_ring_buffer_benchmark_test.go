package ringbuffer

import (
	"strings"
	"testing"
)

func BenchmarkAtomicRingBuffer_Sync(b *testing.B) {
	rb := NewAtomicRingBuffer(1024)
	data := []byte(strings.Repeat("a", 512))
	buf := make([]byte, 512)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Write(data)
		rb.Read(buf)
	}
}

func BenchmarkAtomicRingBuffer_AsyncRead(b *testing.B) {
	rb := NewAtomicRingBuffer(1024)
	data := []byte(strings.Repeat("a", 512))
	buf := make([]byte, 512)

	go func() {
		for {
			rb.Read(buf)
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Write(data)
	}
}

func BenchmarkAtomicRingBuffer_AsyncWrite(b *testing.B) {
	rb := NewAtomicRingBuffer(1024)
	data := []byte(strings.Repeat("a", 512))
	buf := make([]byte, 512)

	go func() {
		for {
			rb.Write(data)
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Read(buf)
	}
}
