package stage

import (
	"testing"
	"time"
)

func TestRingBuf(t *testing.T) {
	done := make(chan struct{})

	in := make(chan int)
	ringOut := RingBuf(done, in, 5)

	for i := 0; i < 3; i++ {
		in <- i
	}
	numsOut := drain(Take(done, ringOut, 3))
	expected := []int{0, 1, 2}
	for i := 0; i < len(expected); i++ {
		if numsOut[i] != expected[i] {
			t.Fatalf("expected=%v but actual=%v", expected, numsOut)
		}
	}

	for i := 0; i < 5; i++ {
		in <- i
	}
	numsOut = drain(Take(done, ringOut, 5))
	expected = []int{0, 1, 2, 3, 4}
	for i := 0; i < len(expected); i++ {
		if numsOut[i] != expected[i] {
			t.Fatalf("expected=%v but actual=%v", expected, numsOut)
		}
	}

	for i := 0; i < 5; i++ {
		in <- i
	}
	for i := 100; i < 103; i++ {
		in <- i
	}
	time.Sleep(5 * time.Millisecond)
	numsOut = drain(Take(done, ringOut, 5))
	expected = []int{3, 4, 100, 101, 102}
	for i := 0; i < len(expected); i++ {
		if numsOut[i] != expected[i] {
			t.Fatalf("expected=%v but actual=%v", expected, numsOut)
		}
	}

	for i := 0; i < 5; i++ {
		in <- i
	}
	for i := 100; i < 108; i++ {
		in <- i
	}
	time.Sleep(5 * time.Millisecond)
	numsOut = drain(Take(done, ringOut, 5))
	expected = []int{103, 104, 105, 106, 107}
	for i := 0; i < len(expected); i++ {
		if numsOut[i] != expected[i] {
			t.Fatalf("expected=%v but actual=%v", expected, numsOut)
		}
	}

	for i := 0; i < 5; i++ {
		in <- i
	}
	for i := 100; i < 103; i++ {
		in <- i
	}
	time.Sleep(5 * time.Millisecond)
	numsOut = drain(Take(done, ringOut, 2))
	expected = []int{3, 4}
	for i := 0; i < len(expected); i++ {
		if numsOut[i] != expected[i] {
			t.Fatalf("expected=%v but actual=%v", expected, numsOut)
		}
	}
	for i := 200; i < 204; i++ {
		in <- i
	}
	time.Sleep(5 * time.Millisecond)
	numsOut = drain(Take(done, ringOut, 5))
	expected = []int{102, 200, 201, 202, 203}
	for i := 0; i < len(expected); i++ {
		if numsOut[i] != expected[i] {
			t.Fatalf("expected=%v but actual=%v", expected, numsOut)
		}
	}
}

func drain[T any](in <-chan T) []T {
	var out []T
	for v := range in {
		out = append(out, v)
	}
	return out
}

func TestRingBuf_Cancelled(t *testing.T) {
	done := make(chan struct{})
	in := make(chan int)

	ringSize := 1000
	ringBuf := RingBuf(done, in, ringSize)

	for i := 0; i < ringSize; i++ {
		in <- i
	}

	go func() {
		defer func() {
			recover()
		}()
		for i := 0; i < 1_000_000_000_000_000_000; i++ {
			in <- i
		}
	}()

	time.Sleep(5 * time.Millisecond)

	done <- struct{}{}

	// just in case there comes a day when that is still running and other tests are running at the same time...
	close(in)

	out := drain(ringBuf)
	for i := 0; i < len(out); i++ {
		if out[i] < ringSize {
			t.Fatalf("The initial [0, %d) should not be in the ring buffer after executing this long, found the value %d", ringSize, out[i])
		}
	}
}
