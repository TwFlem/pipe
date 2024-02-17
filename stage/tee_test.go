package stage

import (
	"math"
	"testing"
	"time"
)

func TestTee(t *testing.T) {
	done := make(chan struct{})
	expectedMin := 10
	expectedMax := 20
	nums := genNumsRange(expectedMin, expectedMax+1)
	in := make(chan int)
	go func() {
		defer close(in)
		for i := range nums {
			in <- nums[i]
		}
	}()

	o1, o2 := Tee(done, in)

	var outNums1 []int
	var outNums2 []int

	o1Done := false
	o2Done := false
	after := time.After(time.Second)
	for {
		if o1Done && o2Done {
			break
		}
		select {
		case <-after:
			t.Fatal("probably deadlocked")
		case v, ok := <-o1:
			if !ok {
				o1Done = true
			} else {
				outNums1 = append(outNums1, v)
			}
		case v, ok := <-o2:
			if !ok {
				o2Done = true
			} else {
				outNums2 = append(outNums2, v)
			}
		}
	}

	expectedLen := expectedMax - expectedMin + 1
	if len(outNums1) != expectedLen {
		t.Fatalf("Expected out1 to have %d nums but had %d\n", expectedLen, len(outNums1))
	}
	if len(outNums2) != expectedLen {
		t.Fatalf("Expected out2 to have %d nums but had %d\n", expectedLen, len(outNums2))
	}

	max1 := math.MinInt
	min1 := math.MaxInt
	for i := 0; i < len(outNums1); i++ {
		if outNums1[i] < min1 {
			min1 = outNums1[i]
		}
		if outNums1[i] > max1 {
			max1 = outNums1[i]
		}
	}
	if min1 != expectedMin {
		t.Fatalf("Expected out1 to have min of %d but had %d\n", expectedMin, min1)
	}
	if max1 != expectedMax {
		t.Fatalf("Expected out1 to have max of %d but had %d\n", expectedMax, max1)
	}

	max2 := math.MinInt
	min2 := math.MaxInt
	for i := 0; i < len(outNums2); i++ {
		if outNums2[i] < min2 {
			min2 = outNums2[i]
		}
		if outNums2[i] > max2 {
			max2 = outNums2[i]
		}
	}
	if min2 != expectedMin {
		t.Fatalf("Expected out2 to have min of %d but had %d\n", expectedMin, min2)
	}
	if max2 != expectedMax {
		t.Fatalf("Expected out2 to have max of %d but had %d\n", expectedMax, max2)
	}
}

func TestTee_Cancelled(t *testing.T) {
	done := make(chan struct{})
	maxVal := 100
	maxLen := maxVal + 1
	nums := genNums(maxLen)
	in := make(chan int)
	go func() {
		defer close(in)
		for i := range nums {
			in <- nums[i]
		}
	}()

	o1, o2 := Tee(done, in)

	var outNums1 []int
	var outNums2 []int

	o1Done := false
	o2Done := false
	after := time.After(time.Second)
	count := 0
	for {
		if o1Done && o2Done {
			break
		}
		select {
		case <-after:
			t.Fatal("probably deadlocked")
		case v, ok := <-o1:
			if !ok {
				o1Done = true
			} else {
				outNums1 = append(outNums1, v)
				count++
				if count == 10 {
					done <- struct{}{}
				}
			}
		case v, ok := <-o2:
			if !ok {
				o2Done = true
			} else {
				outNums2 = append(outNums2, v)
			}
		}
	}

	if len(outNums1) >= maxLen {
		t.Fatalf("Expected out1 to have less than %d but had %d\n", maxLen, len(outNums1))
	}
	if len(outNums2) >= maxLen {
		t.Fatalf("Expected out2 to have less than %d but had %d\n", maxLen, len(outNums2))
	}
}
