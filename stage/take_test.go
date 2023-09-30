package stage

import "testing"

func TestTake(t *testing.T) {
	done := make(chan struct{})

	nums := genNums(100)
	in := make(chan int)

	go func() {
		defer close(in)
		for v := range nums {
			in <- v
		}
	}()

	expectedCount := 10
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

	nums := genNums(100)
	in := make(chan int)

	go func() {
		defer close(in)
		count := 0
		for v := range nums {
			in <- v
			count++
			if count == 50 {
				done <- struct{}{}
			}
		}
	}()

	takeCount := 90
	out := Take(done, in, takeCount)

	count := 0
	for range out {
		count++
	}

	if count >= takeCount {
		t.Fatalf("Expected to receive less than %d but received %d\n", takeCount, count)
	}
}
