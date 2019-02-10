package main

import (
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	source := produce()

	wg.Add(1)
	go func() {
		process(source)
		wg.Done()
	}()
	wg.Wait()
}

func produce() <-chan int {
	res := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			res <- i
		}
		close(res)
	}()
	return res
}

func process(in <-chan int) {
	for i := range in {
		slowCall(i)
	}
}

func slowCall(seqNum int) int {
	time.Sleep(1 * time.Second)
	return seqNum
}
