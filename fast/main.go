package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	source := producer()

	wg.Add(1)
	go func() {
		setup(source)
		wg.Done()
	}()
	wg.Wait()
}

func producer() <-chan int {
	res := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			res <- i
		}
		close(res)
	}()
	return res
}

func setup(in <-chan int) {
	resChan := make(chan int)
	doneChan := make(chan bool)
	wg := sync.WaitGroup{}

	wg.Add(5)
	for i := 0; i < 5; i++ {
		go worker(in, resChan, &wg)
	}
	go processResults(resChan, doneChan)
	wg.Wait()
	doneChan <- true
}

func processResults(in chan int, done <-chan bool) {
	for {
		select {
		case i := <-in:
			fmt.Printf("Response to %d\n", i)
		case <-done:
			for i := range in {
				fmt.Printf("Response to %d\n", i)
			}
			close(in)
			return
		}
	}
}

func worker(in <-chan int, out chan<- int, wg *sync.WaitGroup) {
	for i := range in {
		res := slowCall(i)
		out <- res
	}
	wg.Done()
}

func slowCall(seqNum int) int {
	time.Sleep(time.Second)
	return seqNum
}
