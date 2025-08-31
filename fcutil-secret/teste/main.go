package main

import "fmt"

func main() {
	evento := []string{"teste", "teste2", "teste3", "teste4"}

	evento = append(evento[:0], evento[1:]...)
	fmt.Println(evento)
}
