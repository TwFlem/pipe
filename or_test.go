package pipe

import (
	"testing"
	"time"
)

func TestOr(t *testing.T) {
	done := make(chan struct{})

	first := make(chan int)
	chans := [](<-chan int){first}

	for i := 0; i < 2; i++ {
		in := make(chan int)
		go func(innerIn chan int, v int) {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			<-ticker.C
			innerIn <- v
		}(in, i)
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

	for i := 0; i < 2; i++ {
		in := make(chan int)
		go func(innerIn chan int, v int) {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			<-ticker.C
			innerIn <- v
		}(in, i)
		chans = append(chans, in)
	}

	go func() {
		<-time.After(500 * time.Millisecond)
		first <- 1337
	}()

	go func() {
		done <- struct{}{}
	}()

	if val, ok := <-Or(done, chans...); ok {
		t.Fatalf("Expected or to finish with not ok but got %d\n", val)
	}
}
