package pipe

import "testing"

func TestOrDone(t *testing.T) {
	done := make(chan struct{})
	out := make(chan int)

	stopInfiniteLoop := make(chan struct{})
	go func() {
		defer close(out)
		i := 0
		for {
			select {
			case out <- i:
				i++
			case <-stopInfiniteLoop:
				return
			}
		}
	}()

	var agg []int
	for v := range OrDone(done, out) {
		if len(agg) != 10 {
			agg = append(agg, v)
		} else {
			done <- struct{}{}
		}
	}

	stopInfiniteLoop <- struct{}{}
}
