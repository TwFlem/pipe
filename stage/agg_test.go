package stage

import (
	"math"
	"testing"
)

func TestAgg_HasFullLastChunk(t *testing.T) {
	done := make(chan struct{})
	countTotal := 20
	batchSize := 5
	in := make(chan int)

	out := Agg(done, in, batchSize)
	go func() {
		defer close(in)
		for i := 0; i < countTotal; i++ {
			in <- i
		}
	}()

	var aggOfAgg [][]int
	for v := range out {
		aggOfAgg = append(aggOfAgg, v)
	}

	expectedNumChunks := math.Ceil(float64(countTotal) / float64(batchSize))
	if len(aggOfAgg) != int(expectedNumChunks) {
		t.Fatalf("Expected %d num of chunks, but received %d\n", int(expectedNumChunks), len(aggOfAgg))
	}

	fatalF := func(exp, act int) { t.Fatalf("Expected chunk len to be %d but got %d\n", exp, act) }
	if len(aggOfAgg[0]) != batchSize {
		fatalF(batchSize, len(aggOfAgg[0]))
	}
	if len(aggOfAgg[1]) != batchSize {
		fatalF(batchSize, len(aggOfAgg[1]))
	}
	if len(aggOfAgg[2]) != batchSize {
		fatalF(batchSize, len(aggOfAgg[2]))
	}
	if len(aggOfAgg[3]) != batchSize {
		fatalF(batchSize, len(aggOfAgg[3]))
	}
}

func TestAgg_HasLittleLastChunk(t *testing.T) {
	done := make(chan struct{})
	countTotal := 17
	batchSize := 5
	in := make(chan int)

	out := Agg(done, in, batchSize)
	go func() {
		defer close(in)
		for i := 0; i < countTotal; i++ {
			in <- i
		}
	}()

	var aggOfAgg [][]int
	for v := range out {
		aggOfAgg = append(aggOfAgg, v)
	}

	expectedNumChunks := math.Ceil(float64(countTotal) / float64(batchSize))
	if len(aggOfAgg) != int(expectedNumChunks) {
		t.Fatalf("Expected %d num of chunks, but received %d\n", int(expectedNumChunks), len(aggOfAgg))
	}

	fatalF := func(exp, act int) { t.Fatalf("Expected chunk len to be %d but got %d\n", exp, act) }
	if len(aggOfAgg[0]) != batchSize {
		fatalF(batchSize, len(aggOfAgg[0]))
	}
	if len(aggOfAgg[1]) != batchSize {
		fatalF(batchSize, len(aggOfAgg[1]))
	}
	if len(aggOfAgg[2]) != batchSize {
		fatalF(batchSize, len(aggOfAgg[2]))
	}

	expectedLastChunk := countTotal % batchSize
	if len(aggOfAgg[3]) != expectedLastChunk {
		fatalF(expectedLastChunk, len(aggOfAgg[3]))
	}

}

func TestAgg_Cancelled_NoLastChunk(t *testing.T) {
	done := make(chan struct{})
	countTotal := 12
	batchSize := 3
	stopSendingAfter := 10
	in := make(chan int)

	out := Agg(done, in, batchSize)
	go func() {
		defer close(in)
		for i := 0; i < countTotal; i++ {
			in <- i
			if stopSendingAfter == i {
				done <- struct{}{}
			}
		}
	}()

	var aggOfAgg [][]int
	for v := range out {
		aggOfAgg = append(aggOfAgg, v)
	}

	expectedNumChunks := (countTotal / batchSize) - 1
	if len(aggOfAgg) != expectedNumChunks {
		t.Fatalf("Expected %d num of chunks, but received %d\n", expectedNumChunks, len(aggOfAgg))
	}

	fatalF := func(exp, act int) { t.Fatalf("Expected chunk len to be %d but got %d\n", exp, act) }
	if len(aggOfAgg[0]) != batchSize {
		fatalF(batchSize, len(aggOfAgg[0]))
	}
	if len(aggOfAgg[1]) != batchSize {
		fatalF(batchSize, len(aggOfAgg[1]))
	}
	if len(aggOfAgg[2]) != batchSize {
		fatalF(batchSize, len(aggOfAgg[2]))
	}
}

func TestAgg_FirstChunkFull(t *testing.T) {
	done := make(chan struct{})
	countTotal := 5
	batchSize := 5
	in := make(chan int)

	out := Agg(done, in, batchSize)
	go func() {
		defer close(in)
		for i := 0; i < countTotal; i++ {
			in <- i
		}
	}()

	var aggOfAgg [][]int
	for v := range out {
		aggOfAgg = append(aggOfAgg, v)
	}

	expectedNumChunks := math.Ceil(float64(countTotal) / float64(batchSize))
	if len(aggOfAgg) != int(expectedNumChunks) {
		t.Fatalf("Expected %d num of chunks, but received %d\n", int(expectedNumChunks), len(aggOfAgg))
	}

	fatalF := func(exp, act int) { t.Fatalf("Expected chunk len to be %d but got %d\n", exp, act) }
	if len(aggOfAgg[0]) != batchSize {
		fatalF(batchSize, len(aggOfAgg[0]))
	}
}

func TestAgg_FirstChunkLittle(t *testing.T) {
	done := make(chan struct{})
	countTotal := 3
	batchSize := 5
	in := make(chan int)

	out := Agg(done, in, batchSize)
	go func() {
		defer close(in)
		for i := 0; i < countTotal; i++ {
			in <- i
		}
	}()

	var aggOfAgg [][]int
	for v := range out {
		aggOfAgg = append(aggOfAgg, v)
	}

	expectedNumChunks := math.Ceil(float64(countTotal) / float64(batchSize))
	if len(aggOfAgg) != int(expectedNumChunks) {
		t.Fatalf("Expected %d num of chunks, but received %d\n", int(expectedNumChunks), len(aggOfAgg))
	}

	fatalF := func(exp, act int) { t.Fatalf("Expected chunk len to be %d but got %d\n", exp, act) }
	expectedLastChunk := countTotal % batchSize
	if len(aggOfAgg[0]) != expectedLastChunk {
		fatalF(expectedLastChunk, len(aggOfAgg[0]))
	}
}

func TestAgg_Cancelled_NoFirstChunk(t *testing.T) {
	done := make(chan struct{})
	countTotal := 12
	batchSize := 6
	stopSendingAfter := 2
	in := make(chan int)

	out := Agg(done, in, batchSize)
	go func() {
		defer close(in)
		for i := 0; i < countTotal; i++ {
			in <- i
			if stopSendingAfter == i {
				done <- struct{}{}
			}
		}
	}()

	var aggOfAgg [][]int
	for v := range out {
		aggOfAgg = append(aggOfAgg, v)
	}

	expectedNumChunks := 0
	if len(aggOfAgg) != int(expectedNumChunks) {
		t.Fatalf("Expected %d num of chunks, but received %d\n", int(expectedNumChunks), len(aggOfAgg))
	}

}
