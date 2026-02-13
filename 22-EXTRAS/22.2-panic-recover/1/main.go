package main

import "fmt"

func panic1() {
	panic("panic1")
}

func panic2() {
	panic("panic2")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			if r == "panic1" {
				fmt.Println("Recovered from panic1")
			} else if r == "panic2" {
				fmt.Println("Recovered from panic2")
			}
		}
	}()

	//panic1()
	panic2()
}
