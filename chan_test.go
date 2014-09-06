package infchan

import (
	"runtime"
	"testing"
	"time"
)

func TestLink(t *testing.T) {
	in := make(chan int)
	out := make(chan int)
	Link(in, out)
	n := 102400
	for i := 0; i < n; i++ {
		in <- i
	}
	for i := 0; i < n; i++ {
		<-out
	}
}

func TestKill(t *testing.T) {
	for i := 0; i < 1024; i++ {
		in := make(chan int)
		out := make(chan int)
		kill := Link(in, out)
		close(kill)
	}
	time.Sleep(time.Millisecond * 100)
	if runtime.NumGoroutine() > 32 {
		t.Fatal()
	}
}

func BenchmarkChan(b *testing.B) {
	in := make(chan bool)
	out := make(chan bool)
	Link(in, out)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		in <- true
		<-out
	}
}
