package channels

import (
	"sync/atomic"
	"unsafe"
)

//node represents each node in the unbounded queue.
//value is actual value that has been queued.
//next points to the next element in the queue.
type node struct {
	value interface{}
	next  *node
}

//noPointer is a simple sentinel pointer for initial head and tails.
//This will be used to make CAS operations more sensible.
var noPointer = &node{
	value: nil,
	next:  nil,
}

//UnboundedChannel represents the actual channel and encapsulates its elements.
//Use the special NewUnboundedChannel to create such a channel.
type UnboundedChannel struct {
	head *node
	tail *node
}

//NewUnboundedChannel creates a new unbounded channel with two simple sentinel head and tail pointers.
func NewUnboundedChannel() *UnboundedChannel {
	sentinel := &node{
		value: nil,
		next:  noPointer,
	}

	return &UnboundedChannel{
		head: sentinel,
		tail: sentinel,
	}
}

//Enqueue enqueues the given element in the channel. We perform a series of CAS operations to
//successfully enqueue one element. If CAS fails, we would retry the operations until one success
//link happens.
//at first, when the tail has not been changed and the its next pointer points to nothing special,
//we try to atomically sets its next pointer to our new node. If we successfully managed to do that,
//then we would update the tail pointer to point to our new node.
//if for some reason the next pointer of the tail is not nothing (noPointer), then we would co-operate
//with others and make the tail to point to the actual last element.
func (ch *UnboundedChannel) Enqueue(value interface{}) {
	newNode := &node{
		value: value,
		next:  noPointer,
	}

	for {
		tail := ch.tail
		next := tail.next

		if tail == ch.tail {
			if next == noPointer {
				nextPtr := (*unsafe.Pointer)(unsafe.Pointer(&ch.tail.next))
				if atomic.CompareAndSwapPointer(nextPtr, unsafe.Pointer(next), unsafe.Pointer(newNode)) {
					tailPtr := (*unsafe.Pointer)(unsafe.Pointer(&ch.tail))
					atomic.CompareAndSwapPointer(tailPtr, unsafe.Pointer(tail), unsafe.Pointer(newNode))
					break
				}
			} else {
				tailPtr := (*unsafe.Pointer)(unsafe.Pointer(&ch.tail))
				atomic.CompareAndSwapPointer(tailPtr, unsafe.Pointer(tail), unsafe.Pointer(next))
			}
		}
	}
}

//Dequeue removes one element from the channel. if the channel is empty, then currently it spins until it becomes
//non-empty (improvement point). otherwise, it simply tries to get the first element and then CAS the head pointer
//to its next element. when one goroutine tries to put something into channel and another one is trying to get
//something out of it, we can exchange those values directly and skip the overhead of queuing and de-queuing.
func (ch *UnboundedChannel) Dequeue() interface{} {
	for {
		head := ch.head
		tail := ch.tail
		firstElement := head.next

		if head == ch.head {
			if head == tail {
				if firstElement == noPointer {
					continue
				}
				tailPtr := (*unsafe.Pointer)(unsafe.Pointer(&ch.tail))
				atomic.CompareAndSwapPointer(tailPtr, unsafe.Pointer(tail), unsafe.Pointer(firstElement))
			} else {
				value := firstElement.value
				headPtr := (*unsafe.Pointer)(unsafe.Pointer(&ch.head))
				if atomic.CompareAndSwapPointer(headPtr, unsafe.Pointer(head), unsafe.Pointer(firstElement)) {
					return value
				}
			}
		}
	}
}
