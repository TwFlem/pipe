package stage

import (
	"sync"
	"testing"
)

func TestBoundedBuf(t *testing.T) {
	done := make(chan struct{})
	var dropWg sync.WaitGroup
	onDrop := func() {
		defer dropWg.Done()
	}

	in := make(chan int)
	boundedOut := BoundedBuf(done, in, 5, BoundedBufWithOnDrop(onDrop))

	for i := 0; i < 3; i++ {
		in <- i
	}
	numsOut := drain(Take(done, boundedOut, 3))
	expected := []int{0, 1, 2}
	intSliceEquals(t, numsOut, expected)

	for i := 0; i < 5; i++ {
		in <- i
	}
	numsOut = drain(Take(done, boundedOut, 5))
	expected = []int{0, 1, 2, 3, 4}
	intSliceEquals(t, numsOut, expected)

	for i := 0; i < 5; i++ {
		in <- i
	}
	dropWg.Add(3)
	for i := 100; i < 103; i++ {
		in <- i
	}
	dropWg.Wait()
	numsOut = drain(Take(done, boundedOut, 5))
	expected = []int{0, 1, 2, 3, 4}
	intSliceEquals(t, numsOut, expected)

	for i := 0; i < 5; i++ {
		in <- i
	}
	numsOut = drain(Take(done, boundedOut, 2))
	expected = []int{0, 1}
	intSliceEquals(t, numsOut, expected)
	dropWg.Add(3)
	for i := 100; i < 105; i++ {
		in <- i
	}
	dropWg.Wait()
	numsOut = drain(Take(done, boundedOut, 5))
	expected = []int{2, 3, 4, 100, 101}
	intSliceEquals(t, numsOut, expected)
}

func TestBoundedBuf_Cancelled(t *testing.T) {
	done := make(chan struct{})
	in := make(chan int)

	dropCount := 5
	count := 0
	var mut sync.Mutex
	onDrop := func() {
		mut.Lock()
		count++
		if count == dropCount {
		}
		mut.Unlock()
	}

	boundedSize := 10
	boundedBuf := BoundedBuf(done, in, boundedSize, BoundedBufWithOnDrop(onDrop))

	go func() {
		for i := 0; i < 20; i++ {
			in <- i
		}
		done <- struct{}{}
	}()

	for range boundedBuf {
	}
}
