package stage

import (
	"sync"
	"time"
)

// OrDone syntatic sugar that allows the caller to either block on
// reading in the next input from some specified channel or cancel
// work when signaling done
func OrDone[T any](done <-chan struct{}, in <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}
	}()
	return out
}

// Agg aggregates input into some fixed sized chunks and sends out the chunks.
// Useful for batching.
// WARNING: Do not assume all chunks flowing out from Agg will be of the specified size.
// If Agg stops executing before the current chunk fills up the specified
// size, the last chunk sent out will not be of size. Instead, it will be copied
// to a smaller chunk.
// NOTE: If Agg is cancelled mid way through a chunk, the results aggregated up until
// cancellation will just be thrown out.
func Agg[T any](done <-chan struct{}, in <-chan T, size int) chan []T {
	out := make(chan []T)
	go func() {
		defer close(out)
		chunk := make([]T, size)
		i := 0
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					if i > 0 {
						smallerChunk := make([]T, i)
						for j := 0; j < len(smallerChunk); j++ {
							smallerChunk[j] = chunk[j]
						}
						select {
						case <-done:
						case out <- smallerChunk:
						}
					}
					return
				}
				chunk[i] = v
				i++
				if i >= size {
					select {
					case <-done:
						return
					case out <- chunk:
						chunk = make([]T, size)
						i = 0
					}
				}

			}
		}
	}()
	return out
}

// FanIn take some number of channels and funnel them back into one channel. FanIn does
// not respect ordering of results and will start streaming out values as soon as results
// from the incoming channels arrive.
func FanIn[T any](done <-chan struct{}, ins ...<-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		var wg sync.WaitGroup
		for i := range ins {
			wg.Add(1)
			go func(in <-chan T) {
				defer wg.Done()
				for v := range OrDone(done, in) {
					select {
					case <-done:
						return
					case out <- v:
					}
				}
			}(ins[i])
		}
		wg.Wait()
	}()
	return out
}

// Bridge takes a channel of channels as input and streams out the values of each of those
// channels on a single output. This is similar
// to how FanIn works except that with the channels of channel input, order is implicitly
// maintained.
func Bridge[T any](done <-chan struct{}, in <-chan <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for v := range OrDone(done, in) {
			for v := range OrDone(done, v) {
				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}

	}()
	return out
}

// Throttle only allow values in to values out once every tick of the delay duration
func Throttle[T any](done <-chan struct{}, in <-chan T, delay time.Duration) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		ticker := time.NewTicker(delay)
		for i := range OrDone(done, in) {
			select {
			case <-done:
				return
			case <-ticker.C:
				select {
				case <-done:
					return
				case out <- i:
				}
			}
		}
	}()
	return out
}

// Take consume some number of values from some stream before exiting
func Take[T any](done <-chan struct{}, in <-chan T, n int) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case out <- <-in:
			}
		}
	}()
	return out
}

// Or pattern takes in multiple channels, waits for at least one of them
// to receive a value before continuing execution.
// This is useful if you need to wait for at least value for some unknown
// amount of channels.
// The amount of goroutines creates is X/2 where X is the number of channels
// passed in
func Or[T any](done <-chan struct{}, in ...<-chan T) <-chan T {
	if len(in) == 0 {
		return nil
	}
	if len(in) == 1 {
		return in[0]
	}

	innerDone := make(chan T)
	go func() {
		defer close(innerDone)

		var val T
		if len(in) == 2 {
			select {
			case <-done:
				return
			case v := <-in[0]:
				val = v
			case v := <-in[1]:
				val = v
			}
		} else {
			select {
			case <-done:
				return
			case v := <-in[0]:
				val = v
			case v := <-in[1]:
				val = v
			case v := <-in[2]:
				val = v
			case v := <-Or(done, append(in[3:], innerDone)...):
				val = v
			}
		}

		innerDone <- val
	}()

	return innerDone
}

// Tee like a T pipe, redirect flow of data to two separate outputs.
// this could be useful for post processing of data
//
// WARNING: be mindful of what you tee out. For example, you probably wouldn't want to tee out the same pointer unless you it wouldn't be mutated down stream
//
// WARNING: This implementation of tee also tightly couples both outputs. The
// next value from in cannot be piped through tee's outputs until both
// of tee's previous outputs have been received.
func Tee[T any](done <-chan struct{}, in <-chan T) (<-chan T, <-chan T) {
	out1 := make(chan T)
	out2 := make(chan T)
	go func() {
		defer close(out1)
		defer close(out2)

		for v := range OrDone(done, in) {
			o1 := out1
			o2 := out2
			for i := 0; i < 2; i++ {
				select {
				case <-done:
					return
				case o1 <- v:
					o1 = nil
				case o2 <- v:
					o2 = nil
				}
			}
		}

	}()
	return out1, out2
}

// Buf convenient means of buffering up results between pipeline steps
func Buf[T any](done <-chan struct{}, in <-chan T, size int) <-chan T {
	out := make(chan T, size)
	go func() {
		defer close(out)
		for v := range OrDone(done, in) {
			select {
			case <-done:
				return
			case out <- v:
			}
		}
	}()
	return out
}

// RingBuf Circular fix sized buffer. When the buffer is full and the receiver
// is busy, the oldest value in the buffer will be evicted and the incoming
// value will be queued.
func RingBuf[T any](done <-chan struct{}, in <-chan T, size int) <-chan T {
	out := make(chan T, size)
	go func() {
		defer close(out)
		for v := range OrDone(done, in) {
			select {
			case <-done:
				return
			case out <- v:
			default:
				<-out
				out <- v
			}
		}
	}()
	return out
}
