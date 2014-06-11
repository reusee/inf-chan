package c

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestLink(t *testing.T) {
	in := make(chan int)
	out := make(chan int)
	Link(in, out)
	n := 1024000
	for i := 0; i < n; i++ {
		in <- i
	}
	print("all sent\n")
	for i := 0; i < n; i++ {
		<-out
	}
	print("all received\n")
	runtime.GC()
}

func TestKill(t *testing.T) {
	for i := 0; i < 1024; i++ {
		in := make(chan int)
		out := make(chan int)
		kill := Link(in, out)
		close(kill)
	}
	time.Sleep(time.Millisecond * 100)
	fmt.Printf("%d\n", runtime.NumGoroutine())
}
