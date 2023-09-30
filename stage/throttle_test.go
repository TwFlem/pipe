package stage

import (
	"testing"
	"time"
)

func TestThrottle(t *testing.T) {
	done := make(chan struct{})
	nums := genNums(1000)
	expectedSent := 50
	eps := 10

	in := make(chan int)
	go func() {
		defer close(in)
		after := time.After(time.Millisecond * time.Duration(expectedSent))
		for i := range nums {
			select {
			case <-after:
				return
			case in <- nums[i]:
			}
		}
	}()

	throttledOut := Throttle(done, in, time.Millisecond)

	count := 0
	for range throttledOut {
		count++
	}

	if count > expectedSent+eps || count < expectedSent-eps {
		t.Fatalf("Expected at most %d with a possible error range of +-%d, but received %d\n", expectedSent, eps, count)
	}
}
