# Unbounded Channel
**Note**: Please note that this was meant to be a simple demonstration of the ABA problem. I didn't find the time to complete the blog post but you should avoid using this as a lib or something!

A simple lock-free implementation of an unbounded channel in Go.

## How to Use
Initialize the `UnboundedChannel` instance with the `channels.NewUnboundedChannel()` builder:
```go
import "github.com/alimate/unbounded-channel/channels"

ch := channels.NewUnboundedChannel()
```
Then use the `Enqueue` or `Dequeue` methods:
```go
ch.Enqueue(42)
ch.Enqueue(31)

fmt.Println(ch.Dequeue())
fmt.Println(ch.Dequeue())
```

## Exchanger
Currently, when two goroutines try to put and get something to/from the channel, we go through the official `Enqueue` and
`Dequeue` process. One way to improve this situation is to use an `Exchanger` of some form and just exchange the values
between those goroutines.

## Wait Queue
When the queue is empty, currently, we just spin until at least one element enqueued. One way to improve this to use a 
wait queue or `sync.Cond` to sleep and wakeup.

## Credits
This data structure is basically the same as the one articulated by *Maged M. Michael and Michael L. Scott* in their
*Simple, Fast, and Practical Non-Blocking and Blocking Concurrent Queue Algorithms* paper.

## License
Copyright 2020 alimate

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
