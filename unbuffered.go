package main

import (
	"fmt"
	"time"
)

func main() {
	// Create a buffered channel with a capacity of 5
	ch := make(chan int, 5)

	// Sending 5 values into the buffered channel
	for i := 1; i <= 5; i++ {
		fmt.Printf("Sending %d\n", i)
		ch <- i
	}
	fmt.Println("All values sent to the buffered channel")

	// Delayed consumption
	time.Sleep(2 * time.Second)
	fmt.Println("Starting to consume values from the buffered channel")

	// Receiving values from the channel
	for i := 1; i <= 5; i++ {
		val := <-ch
		fmt.Printf("Received %d\n", val)
	}
}
