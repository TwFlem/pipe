package stage

import (
	"testing"
)

func TestBridge_Finishes(t *testing.T) {
	done := make(chan struct{})
	min := 0
	ranges := [][]int{
		genNumsRange(min, 5),
		genNumsRange(5, 15),
		genNumsRange(15, 35),
		genNumsRange(35, 45),
		genNumsRange(45, 50),
	}

	chans := make(chan (<-chan int), len(ranges))

	expectedCount := 0
	for i := 0; i < len(ranges); i++ {
		for j := 0; j < len(ranges[i]); j++ {
			expectedCount++
		}
	}

	go func() {
		defer close(chans)
		for i := range ranges {
			c := make(chan int)
			go func(in chan int, nums []int) {
				defer close(c)
				for j := range nums {
					in <- nums[j]
				}
			}(c, ranges[i])
			chans <- c
		}
	}()

	out := Bridge(done, chans)

	count := 0
	last := min - 1
	for v := range OrDone(done, out) {
		if v <= last {
			t.Fatalf("Expected monotonically increasing values but %d came before %d\n", last, v)
		}
		last = v
		count++
	}
	if count != expectedCount {
		t.Fatalf("Expected total of %d numbers but found %d\n", expectedCount, count)
	}
}

func TestBridge_Cancelled(t *testing.T) {
	done := make(chan struct{})
	min := 0
	ranges := [][]int{
		genNumsRange(min, 5),
		genNumsRange(5, 15),
		genNumsRange(15, 35),
		genNumsRange(35, 45),
		genNumsRange(45, 50),
	}

	chans := make(chan (<-chan int), len(ranges))

	expectedCount := 0
	for i := 0; i < len(ranges); i++ {
		for j := 0; j < len(ranges[i]); j++ {
			expectedCount++
		}
	}

	go func() {
		defer close(chans)
		for i := range ranges {
			c := make(chan int)
			go func(in chan int, nums []int) {
				defer close(c)
				for j := range nums {
					in <- nums[j]
					if nums[j] == 27 {
						done <- struct{}{}
					}
				}
			}(c, ranges[i])
			chans <- c
		}
	}()

	out := Bridge(done, chans)

	count := 0
	last := min - 1
	for v := range OrDone(done, out) {
		if v <= last {
			t.Fatalf("Expected monotonically increasing values but %d came before %d\n", last, v)
		}
		last = v
		count++
	}
	if count >= expectedCount {
		t.Fatalf("Expected less than of %d numbers but found %d\n", expectedCount, count)
	}
}
