package main

import "fmt"

type RingBuffer struct {
	ch chan string
}

func NewRingBuffer() *RingBuffer {
	return &RingBuffer{
		ch: make(chan string, 1),
	}

}

func (r *RingBuffer) Send(item string) {
	select {
	case r.ch <- item:
	default:
		<-r.ch
		r.ch <- item
	}
}

func (r *RingBuffer) C() <-chan string {
	return r.ch
}

func main() {
	rb := NewRingBuffer()

	rb.Send("some value")
	rb.Send("some other value")

	complete := false
	for !complete {
		select {
		case v := <-rb.C():
			fmt.Println(v)
		default:
			complete = true
		}
	}
}
