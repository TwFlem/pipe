package stage

import (
	"fmt"
	"testing"
)

func TestInterleave(t *testing.T) {
	done := make(chan struct{})
	put := func(in chan int, nums ...int) {
		for _, v := range nums {
			in <- v
		}
	}
	var in1 chan int
	var in2 chan int
	var in3 chan int
	var interleaveOut <-chan int
	reInit := func() {
		in1 = make(chan int)
		in2 = make(chan int)
		in3 = make(chan int)
		interleaveOut = Interleave(done, in1, in2, in3)
	}

	reInit()
	go func() {
		put(in2, 2)
		close(in2)
	}()
	go func() {
		put(in1, 1)
		close(in1)
	}()
	go func() {
		put(in3, 3)
		close(in3)
	}()
	actual := drain(Take(done, interleaveOut, 3))
	fmt.Println(actual)
	expected := []int{1, 2, 3}
	intSliceEquals(t, actual, expected)

	reInit()
	go func() {
		put(in1, 1, 4, 7)
		close(in1)
	}()
	go func() {
		put(in2, 2, 5, 8)
		close(in2)
	}()
	go func() {
		put(in3, 3, 6, 9)
		close(in3)
	}()
	expected = drain(interleaveOut)
	actual = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	intSliceEquals(t, actual, expected)

	reInit()
	go func() {
		put(in1, 1, 4, 6)
		close(in1)
	}()
	go func() {
		put(in2, 2)
		close(in2)
	}()
	go func() {
		put(in3, 3, 5, 7, 8, 9)
		close(in3)
	}()
	actual = drain(interleaveOut)
	expected = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	intSliceEquals(t, actual, expected)
}

func TestInterleave_Cancelled(t *testing.T) {
	done := make(chan struct{})
	put := func(in chan int, nums ...int) {
		for _, v := range nums {
			in <- v
		}
	}
	var in1 chan int
	var in2 chan int
	var in3 chan int
	var interleaveOut <-chan int
	reInit := func() {
		in1 = make(chan int)
		in2 = make(chan int)
		in3 = make(chan int)
		interleaveOut = Interleave(done, in1, in2, in3)
	}

	reInit()
	go func() {
		put(in1, 1)
		put(in2, 2)
		put(in3, 3)
		done <- struct{}{}
	}()
	actual := drain(interleaveOut)
	expected := []int{1, 2, 3}
	intSliceEquals(t, actual, expected)
}
