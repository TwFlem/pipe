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
