package singlechan

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

var bufferSizes = []int{1, 16, 64, 512, 65536}

const (
	finalValue = "done"
)

// Sets up a receiver for the ringbuffer channel
// and signals complete when it receives finalValue
func recv(wgRecv *sync.WaitGroup, ch <-chan string) {
	wgRecv.Add(1)
	go func() {
		for v := range ch {
			if v == finalValue {
				wgRecv.Done()
			}
		}
	}()
}

func TestSendRecv(t *testing.T) {
	const numSends = 50000
	const bufferSize = 1
	const sendConcurrency = 4

	var wgSend sync.WaitGroup
	var wgRecv sync.WaitGroup

	rb := NewRingBuffer(WithBufSize(bufferSize))
	defer rb.Close()

	recv(&wgRecv, rb.C())

	semaphore := make(chan struct{}, sendConcurrency)
	wgSend.Add(numSends)

	for i := range numSends {
		go func(i int) {
			semaphore <- struct{}{}
			defer func() {
				<-semaphore
				wgSend.Done()
			}()
			rb.Send(strconv.Itoa(i))
		}(i)
	}

	wgSend.Wait()
	rb.Send(finalValue)
	wgRecv.Wait()

	// messages received
}

func BenchmarkParallelSendReceive(b *testing.B) {
	for _, bufSize := range bufferSizes {
		b.Run(fmt.Sprintf("BufSize:%d", bufSize), func(b *testing.B) {
			b.StopTimer()
			var wgSend sync.WaitGroup
			var wgRecv sync.WaitGroup
			rb := NewRingBuffer(WithBufSize(bufSize))
			defer rb.Close()

			recv(&wgRecv, rb.C())

			wgSend.Add(b.N)
			b.StartTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					rb.Send(strconv.Itoa(0))
					wgSend.Done()
				}
			})

			wgSend.Wait()
			rb.Send("done")
			wgRecv.Wait()
		})
	}
}

func BenchmarkSendReceive(b *testing.B) {
	for _, bufSize := range bufferSizes {
		b.Run(fmt.Sprintf("BufSize:%d", bufSize), func(b *testing.B) {
			b.StopTimer()
			var wgSend sync.WaitGroup
			var wgRecv sync.WaitGroup
			rb := NewRingBuffer(WithBufSize(bufSize))
			defer rb.Close()

			recv(&wgRecv, rb.C())

			wgSend.Add(b.N)
			b.StartTimer()
			for range b.N {
				rb.Send(strconv.Itoa(0))
				wgSend.Done()
			}

			wgSend.Wait()
			rb.Send("done")
			wgRecv.Wait()
		})
	}
}
