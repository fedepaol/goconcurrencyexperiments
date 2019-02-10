package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type message struct {
	msg    string
	seqNum int
}

type result struct {
	res    string
	seqNum int
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		setup()
		wg.Done()
	}()
	wg.Wait()
}

func setup() {
	source := producer()

	numWorkers := 100
	resChan := make(chan result)
	doneChan := make(chan bool)
	wg := sync.WaitGroup{}

	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker(source, resChan, &wg)
	}
	go processResults(resChan, doneChan)

	wg.Wait()
	doneChan <- true
}

func processResults(in chan result, done <-chan bool) {
	nextToProcess := 0
	parked := make([]result, 0)
	for {
		select {
		case i := <-in: // always park & drain. Could optimize by checking if i is the one
			parked = append(parked, i)
			drainParked(&parked, &nextToProcess)
		case <-done:
			for i := range in {
				parked = append(parked, i)
				drainParked(&parked, &nextToProcess)
			}
			close(in)
			return
		}
	}
}

func producer() <-chan message {
	res := make(chan message)
	go func() {
		for i := 0; i < 1000; i++ {
			res <- message{seqNum: i}
		}
		close(res)
	}()
	return res
}

// Given the id of the next message to be processed, sequentially processes
// all the parked messages from that id
// This could be optimized by sorting / using a circular buffer instead of scrolling the slice
// each time
func drainParked(parked *[]result, next *int) {
	for {
		var found bool
		for i, val := range *parked {
			if val.seqNum == *next {
				fmt.Printf("Processed %d\n", val.seqNum)
				*next++
				*parked = append((*parked)[:i], (*parked)[i+1:]...)
				found = true
				break
			}
		}
		if !found {
			return
		}
	}
}

func worker(in <-chan message, out chan<- result, wg *sync.WaitGroup) {
	for msg := range in {
		res := slowCall(msg)
		out <- res
	}
	wg.Done()
}

func slowCall(in message) result {
	howLong := rand.Intn(100)
	time.Sleep(time.Duration(howLong) * time.Millisecond)
	return result{
		seqNum: in.seqNum,
	}
}
