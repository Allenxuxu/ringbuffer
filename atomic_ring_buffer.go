package ringbuffer

import (
	"errors"
	"sync"
	"unsafe"
)

var ErrAtomicIsEmpty = errors.New("atomicringbuffer is empty")

// AtomicRingBuffer is a circular buffer that implement io.ReaderWriter interface.
type AtomicRingBuffer struct {
	buf     []byte
	size    int
	r       int // next position to read
	w       int // next position to write
	isEmpty bool
	mu      sync.Mutex
}

// New returns a new AtomicRingBuffer whose buffer has the given size.
func NewAtomicRingBuffer(size int) *AtomicRingBuffer {
	return &AtomicRingBuffer{
		buf:     make([]byte, size),
		size:    size,
		isEmpty: true,
	}
}

// Read reads up to len(p) bytes into p. It returns the number of bytes read (0 <= n <= len(p)) and any error encountered. Even if Read returns n < len(p), it may use all of p as scratch space during the call. If some data is available but not len(p) bytes, Read conventionally returns what is available instead of waiting for more.
// When Read encounters an error or end-of-file condition after successfully reading n > 0 bytes, it returns the number of bytes read. It may return the (non-nil) error from the same call or return the error (and n == 0) from a subsequent call.
// Callers should always process the n > 0 bytes returned before considering the error err. Doing so correctly handles I/O errors that happen after reading some bytes and also both of the allowed EOF behaviors.
func (r *AtomicRingBuffer) Read(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.read(p)
}

func (r *AtomicRingBuffer) read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	if r.isEmpty {
		return 0, ErrIsEmpty
	}

	if r.w > r.r {
		n = r.w - r.r
		if n > len(p) {
			n = len(p)
		}
		copy(p, r.buf[r.r:r.r+n])
		r.r = (r.r + n) % r.size

		if r.w == r.r {
			r.isEmpty = true
		}
		return
	}

	n = r.size - r.r + r.w
	if n > len(p) {
		n = len(p)
	}

	if r.r+n <= r.size {
		copy(p, r.buf[r.r:r.r+n])
	} else {
		c1 := r.size - r.r
		copy(p, r.buf[r.r:r.size])
		c2 := n - c1
		copy(p[c1:], r.buf[0:c2])
	}
	r.r = (r.r + n) % r.size

	if r.w == r.r {
		r.isEmpty = true
	}
	return n, err
}

// ReadByte reads and returns the next byte from the input or ErrIsEmpty.
func (r *AtomicRingBuffer) ReadByte() (b byte, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isEmpty {
		return 0, ErrIsEmpty
	}
	b = r.buf[r.r]
	r.r++
	if r.r == r.size {
		r.r = 0
	}

	if r.w == r.r {
		r.isEmpty = true
	}

	return b, err
}

// Write writes len(p) bytes from p to the underlying buf.
// It returns the number of bytes written from p (0 <= n <= len(p)) and any error encountered that caused the write to stop early.
// Write returns a non-nil error if it returns n < len(p).
// Write must not modify the slice data, even temporarily.
func (r *AtomicRingBuffer) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	free := r.free()
	if free < len(p) {
		r.makeSpace(len(p) - free)
	}
	n = len(p)

	if r.w >= r.r {
		c1 := r.size - r.w
		if c1 >= n {
			copy(r.buf[r.w:], p)
			r.w += n
		} else {
			copy(r.buf[r.w:], p[:c1])
			c2 := n - c1
			copy(r.buf[0:], p[c1+1:])
			r.w = c2
		}
	} else {
		copy(r.buf[r.w:], p)
		r.w += n
	}

	if r.w == r.size {
		r.w = 0
	}

	r.isEmpty = false

	return n, err
}

// WriteByte writes one byte into buffer, and returns ErrIsFull if buffer is full.
func (r *AtomicRingBuffer) WriteByte(c byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.free() < 1 {
		r.makeSpace(1)
	}

	r.buf[r.w] = c
	r.w++

	if r.w == r.size {
		r.w = 0
	}

	r.isEmpty = false

	return nil
}

// Length return the length of available read bytes.
func (r *AtomicRingBuffer) Length() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.length()
}

func (r *AtomicRingBuffer) length() int {
	if r.w == r.r {
		if r.isEmpty {
			return 0
		}
		return r.size
	}

	if r.w > r.r {
		return r.w - r.r
	}

	return r.size - r.r + r.w
}

// Capacity returns the size of the underlying buffer.
func (r *AtomicRingBuffer) Capacity() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.size
}

// Free returns the length of available bytes to write.
func (r *AtomicRingBuffer) Free() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.free()
}

func (r *AtomicRingBuffer) free() int {
	if r.w == r.r {
		if r.isEmpty {
			return r.size
		}
		return 0
	}

	if r.w < r.r {
		return r.r - r.w
	}

	return r.size - r.w + r.r
}

// WriteString writes the contents of the string s to buffer, which accepts a slice of bytes.
func (r *AtomicRingBuffer) WriteString(s string) (n int, err error) {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	buf := *(*[]byte)(unsafe.Pointer(&h))
	return r.Write(buf)
}

// Bytes returns all available read bytes. It does not move the read pointer and only copy the available data.
func (r *AtomicRingBuffer) Bytes() []byte {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.w == r.r {
		if !r.isEmpty {
			buf := make([]byte, r.size)
			copy(buf, r.buf)
			return buf
		}
		return nil
	}

	if r.w > r.r {
		buf := make([]byte, r.w-r.r)
		copy(buf, r.buf[r.r:r.w])
		return buf
	}

	n := r.size - r.r + r.w
	buf := make([]byte, n)

	if r.r+n < r.size {
		copy(buf, r.buf[r.r:r.r+n])
	} else {
		c1 := r.size - r.r
		copy(buf, r.buf[r.r:r.size])
		c2 := n - c1
		copy(buf[c1:], r.buf[0:c2])
	}

	return buf
}

// IsFull returns this AtomicRingBuffer is full.
func (r *AtomicRingBuffer) IsFull() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	return !r.isEmpty && r.w == r.r
}

// IsEmpty returns this AtomicRingBuffer is empty.
func (r *AtomicRingBuffer) IsEmpty() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.isEmpty
}

// Reset the read pointer and writer pointer to zero.
func (r *AtomicRingBuffer) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.r = 0
	r.w = 0
	r.isEmpty = true
}

func (r *AtomicRingBuffer) makeSpace(len int) {
	newSize := r.size + len
	newBuf := make([]byte, newSize)
	oldLen := r.length()
	_, _ = r.read(newBuf)

	r.w = oldLen
	r.r = 0
	r.size = newSize
	r.buf = newBuf
}
