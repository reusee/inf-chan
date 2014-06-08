package c

import (
	"net/http"
	_ "net/http/pprof"
	"reflect"
)

func init() {
	go http.ListenAndServe("0.0.0.0:60000", nil)
}

type _Node struct {
	value *reflect.Value
	next  *_Node
}

func Link(in interface{}, out interface{}) {
	tail := new(_Node)
	head := tail
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
				} else { // in chan
					if !ok { // in chan closed
						outValue.Close()
						return
					} else { // received
						tail.value = &recv
						node := new(_Node)
						tail.next = node
						tail = node
					}
				}
			} else {
				recv, ok := inValue.Recv()
				if !ok { // in chan closed
					outValue.Close()
					return
				}
				tail.value = &recv
				node := new(_Node)
				tail.next = node
				tail = node
			}
		}
	}()
}
