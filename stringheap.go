package main

import (
	"container/heap"
	"time"

	"github.com/araddon/dateparse"
)

type datefinder func(string) time.Time

// implements heap.Interface
type stringHeap struct {
	Strs []string   // the data
	Fdr  datefinder // method used for finding the date
	Max  int        // max heap size
}

func (s stringHeap) Len() int {
	return len(s.Strs)
}

func (s stringHeap) Less(i, j int) bool {
	ti := s.Fdr(s.Strs[i])
	tj := s.Fdr(s.Strs[j])
	return ti.Before(tj)
}

func (s stringHeap) Swap(i, j int) {
	s.Strs[i], s.Strs[j] = s.Strs[j], s.Strs[i]
}

func (s *stringHeap) Push(x interface{}) {
	s.Strs = append(s.Strs, x.(string))
}

func (s *stringHeap) Pop() interface{} {
	n := len(s.Strs) - 1
	x := s.Strs[n]
	s.Strs = s.Strs[:n]
	return x
}

// expensive search for a date somewhere in the string
func dateParseRecursive(datestr string) (time.Time, bool) {
	if len(datestr) == 0 {
		return time.Time{}, false
	}
	t, err := dateparse.ParseAny(datestr)
	if err != nil {
		return dateParseRecursive(datestr[1:])
	}
	return t, true
}

// if a date can't be found, assume it's the current date.
// implements datefinder
func dateOrNow(datestr string) time.Time {
	t, ok := dateParseRecursive(datestr)
	if !ok {
		t = time.Now()
	}
	return t
}

// if a date can't be found, assume it's old
// implements datefinder
func dateOrOld(datestr string) time.Time {
	t, _ := dateParseRecursive(datestr)
	return t
}

// strings come in on the input string and they are added to the heap.
// Once the heap reaches its maximum size, it will pop items onto the out chan.
// The size of the heap determines the window for corrrectly sorting dates.
func stringHeapWorker(sh *stringHeap, inCh <-chan string, outCh chan<- string) {
	for str := range inCh {
		if sh.Len() > sh.Max {
			outCh <- heap.Pop(sh).(string)
		}
		heap.Push(sh, str)
	}

	// once the in chan is closed, clear out the contents of the heap
	for sh.Len() > 0 {
		outCh <- heap.Pop(sh).(string)
	}
	close(outCh)
}
