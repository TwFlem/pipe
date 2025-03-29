package pipe

import (
	"testing"
)

func genNums(count int) []int {
	nums := make([]int, count)
	for i := range count {
		nums[i] = i
	}
	return nums
}

func genNumsRange(start, upTo int) []int {
	l := upTo - start
	nums := make([]int, upTo-start)
	for i := range l {
		nums[i] = start + i
	}
	return nums
}

func intSliceEquals(t *testing.T, actual, expected []int) {
	if len(actual) != len(expected) {
		t.Fatalf("expected=%v but actual=%v", expected, actual)
	}
	for i := range expected {
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
