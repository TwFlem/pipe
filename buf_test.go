package pipe

import (
	"testing"
	"time"
)

func TestBuf(t *testing.T) {
	done := make(chan struct{})
	in := make(chan int)
	maxOutputBufferSize := 4

	unBlocked := make(chan struct{})
	onBlockReport := func(dur time.Duration) {
		unBlocked <- struct{}{}
	}

	blocked := make(chan struct{})
	onBlock := func() {
		blocked <- struct{}{}
	}

	bufOut := Buf(done, in, maxOutputBufferSize, bufWithMinimumBlockTimeToReport(0), bufWithOnBlockReport(onBlockReport), bufWithOnBlock(onBlock))

	for i := 0; i < 3; i++ {
		in <- i
	}
	numsOut := drain(Take(done, bufOut, 3))
	expected := []int{0, 1, 2}
	intSliceEquals(t, numsOut, expected)

	for i := 0; i < 4; i++ {
		in <- i
	}
	numsOut = drain(Take(done, bufOut, 4))
	expected = []int{0, 1, 2, 3}
	intSliceEquals(t, numsOut, expected)

	for i := 0; i < maxOutputBufferSize+1; i++ {
		in <- i
	}
	<-blocked
	first := <-bufOut
	intCompare(t, 0, first)
	<-unBlocked
	numsOut = drain(Take(done, bufOut, 4))
	expected = []int{1, 2, 3, 4}
	intSliceEquals(t, numsOut, expected)

	for i := 0; i < maxOutputBufferSize+1; i++ {
		in <- i
	}
	<-blocked
	first = <-bufOut
	intCompare(t, 0, first)
	<-unBlocked
	in <- 5
	<-blocked
	second := <-bufOut
	intCompare(t, 1, second)
	<-unBlocked
	in <- 6
	<-blocked
	third := <-bufOut
	intCompare(t, 2, third)
	<-unBlocked
	numsOut = drain(Take(done, bufOut, 4))
	expected = []int{3, 4, 5, 6}
	intSliceEquals(t, numsOut, expected)
}

func TestBuf_Cancelled(t *testing.T) {
	done := make(chan struct{})
	in := make(chan int)

	unBlocked := make(chan struct{})
	onBlockReport := func(dur time.Duration) {
		unBlocked <- struct{}{}
	}

	blocked := make(chan struct{})
	onBlock := func() {
		blocked <- struct{}{}
	}

	maxSize := 2
	bufOut := Buf(done, in, maxSize, bufWithMinimumBlockTimeToReport(0), bufWithOnBlockReport(onBlockReport), bufWithOnBlock(onBlock))

	go func() {
		for i := 0; i < 7; i++ {
			in <- i
		}
		done <- struct{}{}
	}()

	<-blocked
	first := <-bufOut
	intCompare(t, 0, first)
	<-unBlocked
	<-blocked
	second := <-bufOut
	intCompare(t, 1, second)
	<-unBlocked
	<-blocked
	third := <-bufOut
	intCompare(t, 2, third)
	<-unBlocked

	go func() {
		for range blocked {
		}
	}()
	go func() {
		for range unBlocked {
		}
	}()

	for range bufOut {
	}
}
