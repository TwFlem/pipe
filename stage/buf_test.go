package stage

import (
	"testing"
)

// func TestBuf(t *testing.T) {
// 	done := make(chan struct{})
// 	totalNums := 30
// 	nums := genNums(totalNums)
// 	in := make(chan int)
// 	go func() {
// 		defer close(in)
// 		for i := range nums {
// 			in <- nums[i]
// 		}
// 	}()
//
// 	bufSize := 10
// 	out := Buf(done, in, bufSize)
//
// 	fullBufCount := 0
// 	ticker := time.NewTicker(time.Millisecond)
// 	for {
// 		<-ticker.C
// 		bufLen := len(out)
// 		if bufLen != bufSize {
// 			_, ok := <-out
// 			if ok {
// 				t.Fatalf("Expected full buffer with len %d but found %d\n", bufSize, bufLen)
// 			}
// 			break
// 		}
// 		fullBufCount++
// 		for i := 0; i < bufSize; i++ {
// 			<-out
// 		}
// 	}
//
// 	expectedFullBufCount := totalNums / bufSize
// 	if fullBufCount != expectedFullBufCount {
// 		t.Fatalf("Expected full buffer count %d but found %d\n", expectedFullBufCount, fullBufCount)
// 	}
// }

// func TestBuf_Cancelled(t *testing.T) {
// 	done := make(chan struct{})
// 	totalNums := 1000
// 	nums := genNums(totalNums)
// 	in := make(chan int)
// 	go func() {
// 		defer close(in)
// 		for i := range nums {
// 			in <- nums[i]
// 		}
// 	}()
//
// 	bufSize := 50
// 	out := Buf(done, in, bufSize)
//
// 	maxFullBufCount := 3
// 	fullBufCount := 0
// 	ticker := time.NewTicker(time.Millisecond)
// loop:
// 	for {
// 		select {
// 		case <-done:
// 			break loop
// 		case <-ticker.C:
// 			bufLen := len(out)
// 			if bufLen != bufSize {
// 				_, ok := <-out
// 				if ok {
// 					t.Fatalf("Expected full buffer with len %d but found %d\n", bufSize, bufLen)
// 				}
// 				break loop
// 			}
// 			fullBufCount++
// 			for i := 0; i < bufSize; i++ {
// 				<-out
// 			}
// 			if maxFullBufCount == fullBufCount {
// 				done <- struct{}{}
// 				ticker.Stop()
// 				break loop
// 			}
// 		}
// 	}
//
// 	drainCount := 0
// 	for range out {
// 		<-out
// 		drainCount++
// 		if drainCount == bufSize {
// 			t.Fatalf("Got a full sized buffer after signaling done and attempting to drain the buffers out channel\n")
// 		}
// 	}
//
// 	if fullBufCount != maxFullBufCount {
// 		t.Fatalf("Expected full buffer count %d but found %d\n", maxFullBufCount, fullBufCount)
// 	}
// }

func TestBuf_OnBlockReport(t *testing.T) {
}
