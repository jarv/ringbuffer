package singlechan

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var bufferMaxSizes = []int{1, 64, 512}

func TestSendRecv(t *testing.T) {
	const numSends = 50000

	done := make(chan struct{})
	ringBuffer := NewRingBuffer()
	defer ringBuffer.Close()
	values := make([]string, 0, numSends)

	go func() {
		for v := range ringBuffer.C() {
			values = append(values, v)
			if v == strconv.Itoa(numSends) {
				done <- struct{}{}
			}
		}
	}()

	for i := range numSends {
		ringBuffer.Send(strconv.Itoa(i + 1))
	}
	<-done

	assert.Equal(t, strconv.Itoa(numSends), values[len(values)-1], "Did not receive all values!")
}

func BenchmarkSend(b *testing.B) {
	for _, maxSize := range bufferMaxSizes {
		b.Run(fmt.Sprintf("MaxSize:%d", maxSize), func(b *testing.B) {
			b.StopTimer()
			ringBuffer := NewRingBuffer(WithMaxSize(maxSize))

			b.StartTimer()
			for range b.N {
				ringBuffer.Send("1")
			}
		})
	}
}

func BenchmarkParallelSend(b *testing.B) {
	for _, maxSize := range bufferMaxSizes {
		b.Run(fmt.Sprintf("MaxSize:%d", maxSize), func(b *testing.B) {
			b.StopTimer()
			ringBuffer := NewRingBuffer(WithMaxSize(maxSize))

			b.StartTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					ringBuffer.Send("1")
				}
			})
		})
	}
}

func BenchmarkSendReceive(b *testing.B) {
	for _, maxSize := range bufferMaxSizes {
		b.Run(fmt.Sprintf("MaxSize:%d", maxSize), func(b *testing.B) {
			b.StopTimer()
			ringBuffer := NewRingBuffer(WithMaxSize(maxSize))

			b.StartTimer()
			for range b.N {
				ringBuffer.Send("1")
				<-ringBuffer.C()
			}
		})
	}
}
