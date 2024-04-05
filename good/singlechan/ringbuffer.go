package singlechan

import "sync"

const (
	defaultBufSize = 1
)

type RingBuffer struct {
	ch      chan string
	maxSize int
	mu      sync.Mutex
}

func NewRingBuffer(options ...func(*RingBuffer)) *RingBuffer {
	ringBuffer := &RingBuffer{
		maxSize: defaultBufSize,
	}

	for _, o := range options {
		o(ringBuffer)
	}

	ringBuffer.ch = make(chan string, ringBuffer.maxSize)
	return ringBuffer
}

func WithBufSize(maxSize int) func(*RingBuffer) {
	return func(r *RingBuffer) {
		r.maxSize = maxSize
	}
}

func (r *RingBuffer) Send(item string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case r.ch <- item:
	default:
		select {
		case <-r.ch:
			r.ch <- item
		default:
			// channel is empty
			r.ch <- item
		}
	}
}

func (r *RingBuffer) C() <-chan string {
	return r.ch
}

func (r *RingBuffer) Close() {
	close(r.ch)
}
