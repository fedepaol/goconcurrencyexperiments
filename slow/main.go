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
	source := produce()

	wg.Add(1)
	go func() {
		process(source)
		wg.Done()
	}()
	wg.Wait()
}

func produce() <-chan message {
	res := make(chan message)
	go func() {
		for i := 0; i < 1000; i++ {
			res <- message{seqNum: i}
		}
		close(res)
	}()
	return res
}

func process(in <-chan message) {
	for msg := range in {
		result := slowCall(msg)
		fmt.Printf("Processed %d\n", result.seqNum)
	}
}

func slowCall(msg message) result {
	howLong := rand.Intn(100)
	time.Sleep(time.Duration(howLong) * time.Millisecond)
	return result{
		seqNum: msg.seqNum,
	}
}
