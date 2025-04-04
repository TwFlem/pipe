package pipe

import (
	"math"
	"testing"
)

func TestFanIn_Finishes(t *testing.T) {
	done := make(chan struct{})
	expectedMin := 0
	expectedMax := 10
	ranges := [][]int{
		genNumsRange(expectedMin, 2),
		genNumsRange(2, 7),
		genNumsRange(7, expectedMax+1),
	}

	chans := make([](<-chan int), len(ranges))

	expectedCount := 0
	for i := 0; i < len(ranges); i++ {
		for j := 0; j < len(ranges[i]); j++ {
			expectedCount++
		}
	}

	for i := range ranges {
		c := make(chan int)
		chans[i] = c
		go func(in chan int, nums []int) {
			defer close(c)
			for j := range nums {
				in <- nums[j]
			}
		}(c, ranges[i])
	}

	out := FanIn(done, chans...)

	count := 0
	min := math.MaxInt
	max := math.MinInt
	exists := make(map[int]bool)
	for v := range OrDone(done, out) {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}

		if _, ok := exists[v]; ok {
			t.Fatalf("Expected all unique numbers but found duplicate number %d\n", v)
		} else {
			exists[v] = true
		}
		count++
	}

	if min != expectedMin {
		t.Fatalf("Expected min of %d but found %d\n", expectedMin, min)
	}
	if max != expectedMax {
		t.Fatalf("Expected max of %d but found %d\n", expectedMax, max)
	}
	if count != expectedCount {
		t.Fatalf("Expected total of %d numbers but found %d\n", expectedCount, count)
	}
}

func TestFanIn_Cancelled(t *testing.T) {
	done := make(chan struct{})
	countNums := 100
	ranges := [][]int{
		genNumsRange(0, 20),
		genNumsRange(20, 70),
		genNumsRange(70, countNums+1),
	}

	chans := make([](<-chan int), len(ranges))

	for i := range ranges {
		c := make(chan int)
		chans[i] = c
		go func(in chan int, nums []int) {
			defer close(c)
			for j := 0; j < len(nums); j++ {
				in <- nums[j]
			}
		}(c, ranges[i])
	}

	go func() {
		done <- struct{}{}
	}()

	out := FanIn(done, chans...)

	count := 0
	for range out {
		count++
	}

	if count == countNums {
		t.Fatalf("should have exited before output channel received all values")
	}
}
