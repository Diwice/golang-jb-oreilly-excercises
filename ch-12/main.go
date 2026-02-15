package main

import (
	"fmt"
	"math"
	"runtime"
	"sync"
)

// 1
// Had to write some complicated things because of unbuffered channel
// Part of a Worker Pool pattern
func launchGo(goNum int) {
	var wg sync.WaitGroup
	outerChan := make(chan int)

	for i := 1; i <= goNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range 10 {
				outerChan <- i
			}
		}()
	}

	go func() {
		wg.Wait()
		close(outerChan)
	}()

	var swg sync.WaitGroup
	swg.Add(1)
	go func() {
		defer swg.Done()
	outer:
		for {
			select {
			case val, ok := <-outerChan:
				if !ok {
					outerChan = nil
					break outer
				}
				fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
				fmt.Println(val)
			}
		}
	}()

	swg.Wait()
	fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
	fmt.Println("--- End of first function ---")
}

// 2
func launchGoTwo() {
	chanOne, chanTwo := make(chan int, 10), make(chan int, 10)

	go func() {
		for i := range 10 {
			chanOne <- i
		}
		close(chanOne)
	}()

	go func() {
		for j := range 10 {
			chanTwo <- j
		}
		close(chanTwo)
	}()

	for {
		if chanOne == nil && chanTwo == nil {
			break
		}

		fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())

		select {
		case val, ok := <-chanOne:
			if !ok {
				chanOne = nil
				continue
			}

			fmt.Printf("%d: goroutine one\n", val)
		case val, ok := <-chanTwo:
			if !ok {
				chanTwo = nil
				continue
			}

			fmt.Printf("%d: goroutine two\n", val)
		}
	}

	fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
	fmt.Println("--- End of second function ---")
}

// 3
func mapGen() map[int]float64 {
	newMap := make(map[int]float64, 100000)
	for i := range 100000 {
		newMap[i] = math.Sqrt(float64(i))
	}
	return newMap
}

func main() {
	// 1
	launchGo(2)
	// 2
	launchGoTwo()
	// 3
	once := sync.OnceValue(mapGen)
	for i := 1000; i < 100000; i += 1000 {
		val := once()[i]
		fmt.Printf("Sqrt of %d is %f\n", i, val)
	}
}
