package ringbuffer

import (
	"errors"
	"unsafe"
)

var ErrIsEmpty = errors.New("ringbuffer is empty")

// RingBuffer is a circular buffer that implement io.ReaderWriter interface.
type RingBuffer struct {
	buf     []byte
	size    int
	r       int // next position to read
	w       int // next position to write
	isEmpty bool
}

// New returns a new RingBuffer whose buffer has the given size.
func New(size int) *RingBuffer {
	return &RingBuffer{
		buf:     make([]byte, size),
		size:    size,
		isEmpty: true,
	}
}

// NewWithData 特殊场景使用，RingBuffer 会持有data，不会自己申请内存去拷贝
func NewWithData(data []byte) *RingBuffer {
	return &RingBuffer{
		buf:     data,
		size:    len(data),
		r:       0,
		w:       0,
		isEmpty: false,
	}
}

func (r *RingBuffer) RetrieveAll() {
	r.r = 0
	r.w = 0
	r.isEmpty = true
}

func (r *RingBuffer) Retrieve(len int) {
	if len <= 0 {
		return
	}

	if len < r.Length() {
		r.r = (r.r + len) % r.size

		if r.w == r.r {
			r.isEmpty = true
		}
	} else {
		r.RetrieveAll()
	}
}

func (r *RingBuffer) Peek(len int) (first []byte, end []byte) {
	if r.isEmpty {
		return
	}

	if len <= 0 {
		return
	}

	if r.w > r.r {
		n := r.w - r.r // Length
		if n > len {
			n = len
		}

		first = r.buf[r.r : r.r+n]
		return
	}

	n := r.size - r.r + r.w // Length
	if n > len {
		n = len
	}

	if r.r+n <= r.size {
		first = r.buf[r.r : r.r+n]
	} else {
		c1 := r.size - r.r
		first = r.buf[r.r:r.size]
		c2 := n - c1
		end = r.buf[0:c2]
	}

	return
}

func (r *RingBuffer) PeekAll() (first []byte, end []byte) {
	if r.isEmpty {
		return
	}

	if r.w > r.r {
		n := r.w - r.r // Length
		first = r.buf[r.r : r.r+n]
		return
	}

	n := r.size - r.r + r.w // Length
	if r.r+n <= r.size {
		first = r.buf[r.r : r.r+n]
	} else {
		c1 := r.size - r.r
		first = r.buf[r.r:r.size]
		c2 := n - c1
		end = r.buf[0:c2]
	}

	return
}

// Read reads up to len(p) bytes into p. It returns the number of bytes read (0 <= n <= len(p)) and any error encountered. Even if Read returns n < len(p), it may use all of p as scratch space during the call. If some data is available but not len(p) bytes, Read conventionally returns what is available instead of waiting for more.
// When Read encounters an error or end-of-file condition after successfully reading n > 0 bytes, it returns the number of bytes read. It may return the (non-nil) error from the same call or return the error (and n == 0) from a subsequent call.
// Callers should always process the n > 0 bytes returned before considering the error err. Doing so correctly handles I/O errors that happen after reading some bytes and also both of the allowed EOF behaviors.
func (r *RingBuffer) Read(p []byte) (n int, err error) {
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
func (r *RingBuffer) ReadByte() (b byte, err error) {
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
func (r *RingBuffer) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	free := r.Free()
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
			copy(r.buf[0:], p[c1:])
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
func (r *RingBuffer) WriteByte(c byte) error {
	if r.Free() < 1 {
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
func (r *RingBuffer) Length() int {
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
func (r *RingBuffer) Capacity() int {
	return r.size
}

// Free returns the length of available bytes to write.
func (r *RingBuffer) Free() int {
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
func (r *RingBuffer) WriteString(s string) (n int, err error) {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	buf := *(*[]byte)(unsafe.Pointer(&h))
	return r.Write(buf)
}

// Bytes returns all available read bytes. It does not move the read pointer and only copy the available data.
func (r *RingBuffer) Bytes() []byte {
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

// IsFull returns this ringbuffer is full.
func (r *RingBuffer) IsFull() bool {
	return !r.isEmpty && r.w == r.r
}

// IsEmpty returns this ringbuffer is empty.
func (r *RingBuffer) IsEmpty() bool {
	return r.isEmpty
}

// Reset the read pointer and writer pointer to zero.
func (r *RingBuffer) Reset() {
	r.r = 0
	r.w = 0
	r.isEmpty = true
}

func (r *RingBuffer) makeSpace(len int) {
	newSize := r.size + len
	newBuf := make([]byte, newSize)
	oldLen := r.Length()
	_, _ = r.Read(newBuf)

	r.w = oldLen
	r.r = 0
	r.size = newSize
	r.buf = newBuf
}
