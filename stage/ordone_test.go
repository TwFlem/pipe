package stage

import "testing"

func TestOrDone(t *testing.T) {
	done := make(chan struct{})
	overTheMaxCount := 100
	nums := genEvens(overTheMaxCount)
	maxSend := 50
	out := make(chan int)
	go func() {
		defer close(out)
		for i := range nums {
			if i == maxSend {
				return
			}
			out <- nums[i]
		}
		done <- struct{}{}
	}()

	var agg []int
	for v := range OrDone(done, out) {
		agg = append(agg, v)
	}

	if len(agg) > maxSend {
		t.Fatalf("Expected less than %d number sent out of the total of %d defined\n", maxSend, overTheMaxCount)
	}
}
