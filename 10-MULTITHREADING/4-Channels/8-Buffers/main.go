package main

// Thread 1
func main() {
	ch := make(chan string, 2)
	ch <- "Hello"
	ch <- "World"

	println(<-ch)
	println(<-ch)

}
