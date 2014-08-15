package c

import (
	"reflect"
)

type _Node struct {
	value *reflect.Value
	next  *_Node
}

func Link(in interface{}, out interface{}) chan bool {
	tail := new(_Node)
	head := tail
	kill := make(chan bool)
	go func() {
		inValue := reflect.ValueOf(in)
		outValue := reflect.ValueOf(out)
		cases := []reflect.SelectCase{
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
		emptyCase := reflect.SelectCase{
			Dir: reflect.SelectRecv,
		}
		for {
			if head != tail {
				cases[0] = reflect.SelectCase{
					Dir:  reflect.SelectSend,
					Chan: outValue,
					Send: *head.value,
				}
				i, recv, ok := reflect.Select(cases)
				if i == 0 { // sent
					head = head.next
				} else if i == 1 { // in chan
					if !ok { // in chan closed
						outValue.Close()
						return
					} else { // received
						tail.value = &recv
						node := new(_Node)
						tail.next = node
						tail = node
					}
				} else if i == 2 { // kill
					return
				}
			} else {
				cases[0] = emptyCase
				i, recv, ok := reflect.Select(cases)
				if i == 1 { // in chan
					if !ok { // in chan closed
						outValue.Close()
						return
					} else { // received
						tail.value = &recv
						node := new(_Node)
						tail.next = node
						tail = node
					}
				} else if i == 2 { // kill
					return
				} else if i == 0 {
					panic("impossible")
				}
			}
		}
	}()
	return kill
}
