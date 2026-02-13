package main

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

// Number matches any float or integer types
// Based on [this module's]: golang.org/x/exp/constraints Float and Integer interfaces
type Number interface {
	constraints.Float | constraints.Integer
}

// Add accepts two Number type variables
// Returns a Number which is a summary of two input Number(s)
// [Learn more about addition]: https://www.mathisfun.com/numbers/addition.html
func Add[T Number](x, y T) T {
	return x + y
}

func main() {
	fmt.Println(Add(1, 2))
}
