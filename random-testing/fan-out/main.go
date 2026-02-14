package main

import (
	"fmt"
	"sync"
)

func fanOut(out <-chan int, to []chan<- int) {
outer:
	for {
		for i := range to {
			select {
			case val, ok := <-out:
				if !ok {
					out = nil
					break outer
				}

				to[i] <- val
			}
		}
	}

	for j := range to {
		close(to[j])
	}
}

func main() {
	var wg sync.WaitGroup

	mainChan := make(chan int, 5)
	mainChan <- 1
	mainChan <- 2
	mainChan <- 3
	mainChan <- 4
	mainChan <- 5
	close(mainChan)

	subChan1, subChan2 := make(chan int, 3), make(chan int, 3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		fanOut(mainChan, []chan<- int{subChan1, subChan2})
	}()

	wg.Wait()

	fmt.Println("Channel fill: ", len(subChan1), len(subChan2))

	for i := range subChan1 {
		fmt.Println(i)
	}

	for j := range subChan2 {
		fmt.Println(j)
	}
}
