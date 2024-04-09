package main

import "fmt"

type RingBuffer struct {
	inputChannel  chan string
	outputChannel chan string
}

func NewRingBuffer() *RingBuffer {
	return &RingBuffer{
		inputChannel:  make(chan string),
		outputChannel: make(chan string, 1),
	}
}

// https://tanzu.vmware.com/content/blog/a-channel-based-ring-buffer-in-go
func (r *RingBuffer) Run() {
	for v := range r.inputChannel {
		select {
		case r.outputChannel <- v:
		default:
			<-r.outputChannel
			r.outputChannel <- v
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

func main() {
	rb := NewRingBuffer()
	go rb.Run()

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
