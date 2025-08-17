package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Message struct {
	id  int64
	Msg string
}

// Thread 1
func main() {
	c1 := make(chan Message)
	c2 := make(chan Message)
	var i int64 = 0

	go func() {
		for {
			atomic.AddInt64(&i, 1)
			msg := Message{id: i, Msg: "Hello from RabbitMQ"}
			c1 <- msg
			time.Sleep(time.Second)
		}
	}()
	go func() {
		for {
			time.Sleep(time.Second)
			msg := Message{id: i, Msg: "Hello from Kafka"}
			atomic.AddInt64(&i, 1)
			c2 <- msg
		}
	}()
	for {
		//for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Printf("Received from RabbitMQ: ID=%d, Msg=%s\n", msg1.id, msg1.Msg)
		case msg2 := <-c2:
			fmt.Printf("Received from Kafka: ID=%d, Msg=%s\n", msg2.id, msg2.Msg)
		case <-time.After(time.Second * 3):
			fmt.Println("Timeout: No messages received within 3 seconds")

		}
	}
}
