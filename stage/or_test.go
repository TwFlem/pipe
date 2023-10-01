package stage

import (
	"testing"
	"time"
)

func TestOr(t *testing.T) {
	done := make(chan struct{})

	first := make(chan int)
	chans := [](<-chan int){first}

	for i := 1; i < 10; i++ {
		dur := time.Duration(i) * time.Second
		in := make(chan int)
		go func(innerIn chan int, innerDur time.Duration, v int) {
			<-time.After(dur)
			innerIn <- v
		}(in, dur, i)
		chans = append(chans, in)
	}

	expectedVal := 1337
	go func() {
		first <- expectedVal
	}()

	val := <-Or(done, chans...)

	if val != expectedVal {
		t.Fatalf("Expected first value to finish to be %d but got %d\n", expectedVal, val)
	}
}

func TestOr_Cancelled(t *testing.T) {
	done := make(chan struct{})

	first := make(chan int)
	chans := [](<-chan int){first}

	for i := 1; i < 10; i++ {
		dur := time.Duration(i) * time.Second
		in := make(chan int)
		go func(innerIn chan int, innerDur time.Duration, v int) {
			<-time.After(dur)
			innerIn <- v
		}(in, dur, i)
		chans = append(chans, in)
	}

	go func() {
		<-time.After(500 * time.Millisecond)
		first <- 1337
	}()

	go func() {
		<-time.After(time.Millisecond)
		done <- struct{}{}
	}()

	if val, ok := <-Or(done, chans...); ok {
		t.Fatalf("Expected or to finish with not ok but got %d\n", val)
	}
}