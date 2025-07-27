package main

import (
	"fmt"

	"github.com/katianemiranda/goexpert/Packagind/1/math.go"
)

func main() {
	//	m := math.Math{A: 3, B: 2}
	m := math.NewMath(1, 2)
	fmt.Println(m.Add())
}
