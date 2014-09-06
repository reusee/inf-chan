package infchan

import (
	"reflect"
	"sync"
)

type _Node struct {
	value reflect.Value
	next  *_Node
}

var pool = sync.Pool{
	New: func() interface{} {
		return new(_Node)
	},
}

func Link(in interface{}, out interface{}) chan bool {
	tail := new(_Node)
	head := tail
	kill := make(chan bool)
	var p *_Node

	go func() {
		inValue := reflect.ValueOf(in)
		outValue := reflect.ValueOf(out)

		casesA := []reflect.SelectCase{
			{ // send
				Dir:  reflect.SelectSend,
				Chan: outValue,
			},
			{ // receive
				Dir:  reflect.SelectRecv,
				Chan: inValue,
			},
			{ // kill
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(kill),
			},
		}

		casesB := []reflect.SelectCase{
			{ // receive
				Dir:  reflect.SelectRecv,
				Chan: inValue,
			},
			{ // kill
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(kill),
			},
		}

		for {
			if head != tail {
				casesA[0].Send = head.value
				i, recv, ok := reflect.Select(casesA)
				if i == 0 { // sent
					p = head
					head = head.next
					pool.Put(p)
				} else if i == 1 { // in chan
					if !ok { // in chan closed
						outValue.Close()
						return
					} else { // received
						tail.value = recv
						node := pool.Get().(*_Node)
						tail.next = node
						tail = node
					}
				} else if i == 2 { // kill
					return
				}

			} else {
				i, recv, ok := reflect.Select(casesB)
				if i == 0 { // in chan
					if !ok { // in chan closed
						outValue.Close()
						return
					} else { // received
						tail.value = recv
						node := pool.Get().(*_Node)
						tail.next = node
						tail = node
					}
				} else if i == 1 { // kill
					return
				}
			}
		}
	}()

	return kill
}
