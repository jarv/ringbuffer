package dualchan

const (
	defaultBufSize = 1
)

type RingBuffer struct {
	inputChannel  chan string
	outputChannel chan string
	maxSize       int
}

func NewRingBuffer(options ...func(*RingBuffer)) *RingBuffer {
	ringBuffer := &RingBuffer{
		maxSize:      defaultBufSize,
		inputChannel: make(chan string),
	}

	for _, o := range options {
		o(ringBuffer)
	}

	ringBuffer.outputChannel = make(chan string, ringBuffer.maxSize)
	return ringBuffer
}

func WithBufSize(maxSize int) func(*RingBuffer) {
	return func(r *RingBuffer) {
		r.maxSize = maxSize
	}
}

// https://tanzu.vmware.com/content/blog/a-channel-based-ring-buffer-in-go
func (r *RingBuffer) Run() {
	for v := range r.inputChannel {
		select {
		case r.outputChannel <- v:
		default:
			select {
			case <-r.outputChannel:
				r.outputChannel <- v
			default:
				// channel is empty
				r.outputChannel <- v
			}
		}
	}
	close(r.outputChannel)
}

func (r *RingBuffer) C() <-chan string {
	return r.outputChannel
}

func (r *RingBuffer) Send(item string) {
	r.inputChannel <- item
}

func (r *RingBuffer) Close() {
	close(r.outputChannel)
}
