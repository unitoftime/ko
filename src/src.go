package src

import (
	"fmt"
	"math"
)

func add(a, b int) int {
	// return a + b
	fmt.Println("ADD")
	return a + b + int(math.Floor(1.0))
}

var x int
var y = 10
const z = 500.0

type something struct {
	a int
	b float64
}

var (
	a int = 5
	b = 6
)

const (
	c = float32(7)
	d = 8
)
