package main

import "fmt"

func fanIn(in []<-chan int, res chan int) {
	for i := range in {
	inner:
		for {
			select {
			case val, ok := <-in[i]:
				if !ok {
					in[i] = nil
					break inner
				}

				res <- val
			}
		}
	}

	close(res)
}

func main() {
	mockChan1 := make(chan int, 2)
	mockChan1 <- 1
	mockChan1 <- 2
	close(mockChan1)

	mockChan2 := make(chan int, 2)
	mockChan2 <- 3
	mockChan2 <- 4
	close(mockChan2)

	newChan := make(chan int)
	go fanIn([]<-chan int{mockChan1, mockChan2}, newChan)
	for i := range newChan {
		fmt.Println(i)
	}
}