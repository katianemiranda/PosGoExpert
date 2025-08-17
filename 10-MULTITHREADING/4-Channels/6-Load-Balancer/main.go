package main

import (
	"fmt"
	"time"
)

func worker(workerId int, data chan int) {
	for x := range data {
		fmt.Printf("Worker %d received %d\n", workerId, x)
		time.Sleep(time.Second)
	}
}

// Thread 1
func main() {
	data := make(chan int)
	QtdWorkers := 10

	for i := 0; i < QtdWorkers; i++ {
		go worker(i, data)
	}

	for i := 0; i < 100; i++ {
		data <- i
	}
}
