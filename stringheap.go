package main

import (
	"container/heap"
)

// implements heap.Interface
type stringHeap struct {
	Strs []datedString // the data
	Max  int           // max heap size
}

func (s stringHeap) Len() int {
	return len(s.Strs)
}

func (s stringHeap) Less(i, j int) bool {
	return s.Strs[i].Date.Before(s.Strs[j].Date)
}

func (s stringHeap) Swap(i, j int) {
	s.Strs[i], s.Strs[j] = s.Strs[j], s.Strs[i]
}

func (s *stringHeap) Push(x interface{}) {
	(*s).Strs = append((*s).Strs, x.(datedString))
}

func (s *stringHeap) Pop() interface{} {
	n := len((*s).Strs) - 1
	x := (*s).Strs[n]
	(*s).Strs = (*s).Strs[:n]
	return x
}

// strings come in on the input string and they are added to the heap.
// Once the heap reaches its maximum size, it will pop items onto the out chan.
// The size of the heap determines the window for corrrectly sorting dates.
func stringHeapWorker(sh *stringHeap, inCh <-chan datedString, outCh chan<- datedString) {
	for str := range inCh {
		if sh.Len() > sh.Max {
			outCh <- heap.Pop(sh).(datedString)
		}
		heap.Push(sh, str)
	}

	// once the in chan is closed, clear out the contents of the heap
	for sh.Len() > 0 {
		outCh <- heap.Pop(sh).(datedString)
	}
	close(outCh)
}
