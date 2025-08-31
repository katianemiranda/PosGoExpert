package main

import (
	"fmt"

	"github.com/katianemiranda/fcutils-secret/pkg/events"
)

func main() {
	ed := events.NewEventDispatcher()
	fmt.Println(ed)

}
