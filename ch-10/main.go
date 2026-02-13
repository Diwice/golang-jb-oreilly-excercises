package main

import "fmt"
// Add accepts two integers
// Returns int which is a summary of two input integers
// [Learn more about addition]: https://www.mathisfun.com/numbers/addition.html
func Add(x, y int) int {
	return x + y
}

func main() {
	fmt.Println(Add(1, 2))
}
