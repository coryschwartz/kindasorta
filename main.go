package main

import (
	"flag"
	"sync"
)

func main() {
	window := flag.Int("win", 1000, "window-size -- the number of lines to consider while doing the sort")
	showtime := flag.Bool("showtime", false, "print the detected timestamp")
	showsource := flag.Bool("showsource", false, "print the source where the line was generated")

	flag.Parse()

	inCh := make(chan datedString, 0)
	outCh := make(chan datedString, 0)
	outDone := make(chan bool, 0)

	// start the input workers
	wg := &sync.WaitGroup{}
	for _, fn := range flag.Args() {
		wg.Add(1)
		go fileReaderWorker(wg, fn, inCh)
	}

	// heap sort
	sh := &stringHeap{
		Strs: make([]datedString, 0),
		Max:  *window,
	}
	go stringHeapWorker(sh, inCh, outCh)

	// write the output
	go stdoutWorker(outCh, outDone, *showtime, *showsource)

	wg.Wait() // all readers done reading
	close(inCh)
	<-outDone // writer done writing
}
