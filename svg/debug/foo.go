package main

import (
	"fmt"
	"strconv"
	"time"
)

func startShit() {
	go func() {
		defer func() {
			fmt.Println("shutting down")
		}()
		fmt.Println("starting")
		for i := 0; i < 100; i++ {
			fmt.Println("i is " + strconv.Itoa(i))
			time.Sleep(300 * time.Millisecond)
		}
	}()
	time.Sleep(3 * time.Second)
}

func main() {
	startShit()
	for i := 0; i< 100; i++ {
		fmt.Println("still here")
		time.Sleep(1 * time.Second)
	}
}

