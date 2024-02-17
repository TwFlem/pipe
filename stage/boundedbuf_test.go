package stage

import (
	"testing"
	"time"
)

func TestBoundedBuf(t *testing.T) {
	done := make(chan struct{})

	in := make(chan int)
	boundedOut := BoundedBuf(done, in, 5)

	for i := 0; i < 3; i++ {
		in <- i
	}
	numsOut := drain(Take(done, boundedOut, 3))
	expected := []int{0, 1, 2}
	for i := 0; i < len(expected); i++ {
		if numsOut[i] != expected[i] {
			t.Fatalf("expected=%v but actual=%v", expected, numsOut)
		}
	}

	for i := 0; i < 5; i++ {
		in <- i
	}
	numsOut = drain(Take(done, boundedOut, 5))
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
	numsOut = drain(Take(done, boundedOut, 5))
	expected = []int{0, 1, 2, 3, 4}
	for i := 0; i < len(expected); i++ {
		if numsOut[i] != expected[i] {
			t.Fatalf("expected=%v but actual=%v", expected, numsOut)
		}
	}

	for i := 0; i < 5; i++ {
		in <- i
	}
	numsOut = drain(Take(done, boundedOut, 2))
	expected = []int{0, 1}
	for i := 0; i < len(expected); i++ {
		if numsOut[i] != expected[i] {
			t.Fatalf("expected=%v but actual=%v", expected, numsOut)
		}
	}
	for i := 100; i < 105; i++ {
		in <- i
	}
	time.Sleep(5 * time.Millisecond)
	expected = []int{2, 3, 4, 100, 101}
	for i := 0; i < len(expected); i++ {
		if numsOut[i] != expected[i] {
			t.Fatalf("expected=%v but actual=%v", expected, numsOut)
		}
	}
}

func TestBoundedBuf_Cancelled(t *testing.T) {
	done := make(chan struct{})
	in := make(chan int)

	boundedSize := 1000
	boundedBuf := BoundedBuf(done, in, boundedSize)

	for i := 0; i < boundedSize; i++ {
		in <- i
	}
	time.Sleep(5 * time.Millisecond)

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
	time.Sleep(5 * time.Millisecond)

	out := drain(boundedBuf)
	for i := 0; i < len(out); i++ {
		if out[i] != i {
			t.Fatalf("The initial [0, %d) be the only values left in the queue, found the value %d", boundedSize, out[i])
		}
	}
}
