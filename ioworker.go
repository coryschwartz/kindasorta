package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type datedString struct {
	Date   time.Time
	Source string
	Str    string
}

// input workers
func readerWorker(wg *sync.WaitGroup, rdr io.Reader, source string, inCh chan<- datedString) {
	defer wg.Done()

	scan := bufio.NewScanner(rdr)
	for scan.Scan() {
		txt := scan.Text()
		t := findDate(txt)
		inCh <- datedString{t, source, txt}
	}
}

func fileReaderWorker(wg *sync.WaitGroup, filename string, inCh chan<- datedString) {
	f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)

	defer f.Close()
	if err != nil {
		wg.Done()
		return
	}
	readerWorker(wg, f, filename, inCh)
}

func stdinWorker(wg *sync.WaitGroup, inCh chan<- datedString) {
	readerWorker(wg, os.Stdin, "stdin", inCh)
}

// output workers
func stdoutWorker(outCh <-chan datedString, done chan bool, showtime bool, showsource bool) {
	for str := range outCh {
		var s string
		var t time.Time
		if showtime {
			t = str.Date
		}
		if showsource {
			s = str.Source
		}
		fmt.Printf("%s %v %s\n", s, t, str.Str)
	}
	close(done)
}
