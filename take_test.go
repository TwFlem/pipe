package pipe

import "testing"

func TestTake(t *testing.T) {
	done := make(chan struct{})

	nums := genNums(10)
	in := make(chan int)

	go func() {
		defer close(in)
		for v := range nums {
			in <- v
		}
	}()

	expectedCount := 5
	out := Take(done, in, expectedCount)

	count := 0
	for range out {
		count++
	}

	if count != expectedCount {
		t.Fatalf("Expected to receive %d nums but got count %d\n", expectedCount, count)
	}
}

func TestTake_Cancelled(t *testing.T) {
	done := make(chan struct{})
	nums := genNums(5)

	in := make(chan int, len(nums))
	outTake := Take(done, in, 10)

	go func() {
		for _, v := range nums {
			in <- v
		}
		done <- struct{}{}
	}()

	for range outTake {
	}
}
