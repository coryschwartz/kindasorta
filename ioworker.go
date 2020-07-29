package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
)

// input workers
func readerWorker(wg *sync.WaitGroup, rdr io.Reader, inCh chan<- string) {
	defer wg.Done()

	scan := bufio.NewScanner(rdr)
	for scan.Scan() {
		inCh <- scan.Text()
	}
}

func fileReaderWorker(wg *sync.WaitGroup, filename string, inCh chan<- string) {
	f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)

	defer f.Close()
	if err != nil {
		wg.Done()
		return
	}
	readerWorker(wg, f, inCh)
}

func stdinWorker(wg *sync.WaitGroup, inCh chan<- string) {
	readerWorker(wg, os.Stdin, inCh)
}

// output workers
func stdoutWorker(outCh <-chan string, done chan bool) {
	for line := range outCh {
		fmt.Println(line)
	}
	close(done)
}
