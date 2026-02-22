package main

import (
	"fmt"
	dp "github.com/Diwice/golang-jb-oreilly-excercises/ch-15/simplewebapp/internal/data_processor"
	"net/http"
	"os"
)

func main() {
	// set everything up
	ch1 := make(chan []byte, 100)
	ch2 := make(chan dp.Result, 100)
	go dp.DataProcessor(ch1, ch2)
	f, err := os.Create("results.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	go dp.WriteData(ch2, f)
	err = http.ListenAndServe(":8080", dp.NewController(ch1))
	if err != nil {
		fmt.Println(err)
	}
}
