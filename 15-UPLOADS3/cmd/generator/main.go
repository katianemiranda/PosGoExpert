package main

import (
	"fmt"
	"os"
)

func main() {
	i := 0
	for i < 20 {
		f, err := os.Create(fmt.Sprintf("./tmp/file%d.txt", i))
		if err != nil {
			panic(err)
		}
		defer f.Close()
		i++
		f.WriteString("Hello, World!")
	}
}
