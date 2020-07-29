package main

import (
	"flag"
	"sync"
)

func main() {
	heapSize := flag.Int("-max", 1000, "maximum lines of text to keep in memory while sorting")
	flag.Parse()

	inCh := make(chan string, 0)
	outCh := make(chan string, 0)
	outDone := make(chan bool, 0)

	// start the input workers
	wg := &sync.WaitGroup{}
	for _, fn := range flag.Args() {
		wg.Add(1)
		go fileReaderWorker(wg, fn, inCh)
	}

	// heap sort
	sh := &stringHeap{
		Strs: make([]string, 0),
		Fdr:  dateOrOld,
		Max:  *heapSize,
	}
	go stringHeapWorker(sh, inCh, outCh)

	// write the output
	go stdoutWorker(outCh, outDone)

	wg.Wait() // all readers done reading
	close(inCh)
	<-outDone // writer done writing
}
