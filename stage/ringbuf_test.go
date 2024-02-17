package stage

import (
	"sync"
	"testing"
)

func TestRingBuf(t *testing.T) {
	done := make(chan struct{})
	var dropWg sync.WaitGroup
	onDrop := func() {
		defer dropWg.Done()
	}

	in := make(chan int)
	ringOut := RingBuf(done, in, 5, RingBufWithOnDrop(onDrop))

	for i := 0; i < 3; i++ {
		in <- i
	}
	numsOut := drain(Take(done, ringOut, 3))
	expected := []int{0, 1, 2}
	intSliceEquals(t, numsOut, expected)

	for i := 0; i < 5; i++ {
		in <- i
	}
	numsOut = drain(Take(done, ringOut, 5))
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
	numsOut = drain(Take(done, ringOut, 5))
	expected = []int{3, 4, 100, 101, 102}
	intSliceEquals(t, numsOut, expected)

	for i := 0; i < 5; i++ {
		in <- i
	}
	dropWg.Add(8)
	for i := 100; i < 108; i++ {
		in <- i
	}
	dropWg.Wait()
	numsOut = drain(Take(done, ringOut, 5))
	expected = []int{103, 104, 105, 106, 107}
	intSliceEquals(t, numsOut, expected)

	for i := 0; i < 5; i++ {
		in <- i
	}
	dropWg.Add(3)
	for i := 100; i < 103; i++ {
		in <- i
	}
	dropWg.Wait()
	numsOut = drain(Take(done, ringOut, 2))
	expected = []int{3, 4}
	intSliceEquals(t, numsOut, expected)
	dropWg.Add(2)
	for i := 200; i < 204; i++ {
		in <- i
	}
	dropWg.Wait()
	numsOut = drain(Take(done, ringOut, 5))
	expected = []int{102, 200, 201, 202, 203}
	intSliceEquals(t, numsOut, expected)
}

// func TestRingBuf_Cancelled(t *testing.T) {
// 	done := make(chan struct{})
// 	in := make(chan int)
//
// 	ringSize := 10
// 	ringBuf := RingBuf(done, in, ringSize)
//
// 	for i := 0; i < ringSize; i++ {
// 		in <- i
// 	}
//
// 	go func() {
// 		defer func() {
// 			recover()
// 		}()
// 		for i := 0; i < 1_000_000_000_000_000_000; i++ {
// 			in <- i
// 		}
// 	}()
//
// 	time.Sleep(5 * time.Millisecond)
//
// 	done <- struct{}{}
//
// 	// just in case there comes a day when that is still running and other tests are running at the same time...
// 	close(in)
//
// 	out := drain(ringBuf)
// 	for i := 0; i < len(out); i++ {
// 		if out[i] < ringSize {
// 			t.Fatalf("The initial [0, %d) should not be in the ring buffer after executing this long, found the value %d", ringSize, out[i])
// 		}
// 	}
// }
