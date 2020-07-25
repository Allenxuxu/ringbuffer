package LinkBuffer

import "github.com/Allenxuxu/ringbuffer"

type LinkBuffer struct {
	buf  *ringbuffer.RingBuffer
	next *ringbuffer.RingBuffer
}

func New(size int) *LinkBuffer {
	return &LinkBuffer{
		buf:  ringbuffer.New(size),
		next: nil,
	}
}

func (l *LinkBuffer) Write(p []byte) (int, error) {
	length := len(p)
	free := l.buf.Capacity() - l.buf.Length()
	if free < length {
		if l.next == nil {
			l.next = ringbuffer.New(l.buf.Capacity() * 2)
			return l.next.Write(p)
		} else {
			return l.next.Write(p)
		}
	} else {
		return l.buf.Write(p)
	}
}

func (l *LinkBuffer) Read(p []byte) (int, error) {
	length := len(p)
	if l.buf.Length() < length {
		if l.next == nil {
			return l.buf.Read(p)
		} else {
			n, err := l.buf.Read(p)
			if err != nil {
				return n, err
			}

			nextN, err := l.next.Read(p[n:])
			return n + nextN, err
		}
	} else {
		return l.buf.Read(p)
	}
}
