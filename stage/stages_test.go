package stage

import (
	"math"
	"testing"
)

func genNums(count int) []int {
	nums := make([]int, count)
	for i := 0; i < count; i++ {
		nums[i] = i
	}
	return nums
}

func genNumsRange(start, upTo int) []int {
	l := upTo - start
	nums := make([]int, upTo-start)
	for i := 0; i < l; i++ {
		nums[i] = start + i
	}
	return nums
}

func genEvens(count int) []int {
	nums := make([]int, count)
	for i := 0; i < count; i++ {
		nums[i] = 2 * i
	}
	return nums
}

func genOdds(count int) []int {
	nums := make([]int, count)
	for i := 1; i < count; i++ {
		nums[i] = (2 * i) - 1
	}
	return nums
}

func countOdds(nums []int) int {
	count := 0
	for i := range nums {
		if nums[i]%2 == int(math.Abs(1)) {
			count++
		}
	}
	return count
}

func countEvens(nums []int) int {
	count := 0
	for i := range nums {
		if nums[i]%2 == 0 {
			count++
		}
	}
	return count
}

func intSliceEquals(t *testing.T, actual, expected []int) {
	for i := 0; i < len(expected); i++ {
		if actual[i] != expected[i] {
			t.Fatalf("expected=%v but actual=%v", expected, actual)
		}
	}
}

func drain[T any](in <-chan T) []T {
	var out []T
	for v := range in {
		out = append(out, v)
	}
	return out
}

func intCompare(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Fatalf("expected first to be %d, but got %d", expected, actual)
	}
}
