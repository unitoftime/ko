package src

import (
	"fmt"
	"math"
)

func add(a, b int) int {
	fmt.Println("ADD")
	return a + b + int(math.Floor(1.0))
}
