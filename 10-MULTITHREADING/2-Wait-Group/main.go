package main

import (
	"fmt"
	"sync"
	"time"
)

func task(name string, wg *sync.WaitGroup) {
	for i := 0; i < 10; i++ {
		fmt.Printf("%d: Task %s is running\n", i, name)
		time.Sleep(1 * time.Second)
		wg.Done() // Signal that this task is done
	}

}

func main() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(25)
	go task("A", &waitGroup)
	go task("B", &waitGroup)

	go func() {
		for i := 0; i < 5; i++ {
			fmt.Printf("%d: Anonymous task is running\n", i)
			time.Sleep(1 * time.Second)
			waitGroup.Done() // Signal that this task is done
		}
	}()

	waitGroup.Wait()
}
