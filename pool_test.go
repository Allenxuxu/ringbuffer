package ringbuffer

import (
	"testing"
	"time"
)

func BenchmarkTestPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := GetFromPool()
			buf.Length()
			PutInPool(buf)
		}
	})
}

func TestPool(t *testing.T) {
	stop := time.After(time.Second * 3)
	for i := 0; i < 100; i++ {
		go func() {
			for {
				select {
				case <-stop:
					return
				default:
					buf := GetFromPool()
					buf.Length()
					PutInPool(buf)
				}
			}
		}()
	}
	<-stop
}
