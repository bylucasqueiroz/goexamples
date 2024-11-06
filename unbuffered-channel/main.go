package main

import (
	"fmt"
	"time"
)

func main() {
	// Create an unbuffered channel
	ch := make(chan int)

	// Start a goroutine to receive values
	go func() {
		for i := 1; i <= 5; i++ {
			val := <-ch
			fmt.Printf("Received %d\n", val)
			time.Sleep(1 * time.Second) // Simulate some work on each received value
		}
	}()

	// Sending 5 values into the unbuffered channel
	for i := 1; i <= 5; i++ {
		fmt.Printf("Sending %d\n", i)
		ch <- i
		fmt.Printf("Sent %d\n", i)
	}
}
