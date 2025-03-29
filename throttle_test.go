package pipe

import (
	"testing"
	"time"
)

func TestThrottle(t *testing.T) {
	longTestCase(t)

	done := make(chan struct{})
	nums := genNums(100)
	expectedSent := 10
	eps := 3

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
