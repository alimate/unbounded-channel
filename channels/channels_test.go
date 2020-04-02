package channels

import "testing"

func BenchmarkBuiltinChannel(b *testing.B) {
	ch := make(chan int, 100)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ch <- 1
			<-ch
		}
	})
}

func BenchmarkUnboundedChannel(b *testing.B) {
	ch := NewUnboundedChannel()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ch.Enqueue(1)
			ch.Dequeue()
		}
	})
}
