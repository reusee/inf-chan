package c

import (
	"runtime"
	"testing"
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
