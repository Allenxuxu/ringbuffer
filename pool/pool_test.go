package pool

import (
	"github.com/Allenxuxu/ringbuffer"
	"testing"
)

func TestRingBufferPool(t *testing.T) {
	pool := New(1024)

	r := pool.Get()
	if r.Capacity() != 1024 {
		t.Fatal()
	}
	if r.Length() != 0 {
		t.Fatal()
	}
	_, _ = r.Write([]byte("1234"))
	pool.Put(r)

	rr := pool.Get()
	if rr.Capacity() != 1024 {
		t.Fatal()
	}
	if rr.Length() != 4 {
		t.Fatal()
	}

	pool.Put(ringbuffer.New(10))
	rrr := pool.Get()
	if rrr.Capacity() != 10 {
		t.Fatal()
	}
	if rrr.Length() != 0 {
		t.Fatal()
	}
}

func TestDefaultPool(t *testing.T) {
	r := Get()
	if r.Capacity() != 1024 {
		t.Fatal()
	}
	if r.Length() != 0 {
		t.Fatal()
	}
	_, _ = r.Write([]byte("1234"))
	Put(r)

	rr := Get()
	if rr.Capacity() != 1024 {
		t.Fatal()
	}
	if rr.Length() != 4 {
		t.Fatal()
	}

	Put(ringbuffer.New(10))
	rrr := Get()
	if rrr.Capacity() != 10 {
		t.Fatal()
	}
	if rrr.Length() != 0 {
		t.Fatal()
	}
}
