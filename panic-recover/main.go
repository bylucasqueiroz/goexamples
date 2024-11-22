package main

import "fmt"

func main() {
	for i := 0; i < 10; i++ {
		forceFirstPanic(i)
	}

	for i := 0; i < 5; i++ {
		forceSecondPanic(i)
	}
}

func forceFirstPanic(i int) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	if i == 5 {
		panic(fmt.Sprintf("panic on %v", i))
	}

	fmt.Printf("first counting %v\n", i)
}

func forceSecondPanic(i int) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	if i == 4 {
		panic(fmt.Sprintf("second panic on %v", i))
	}

	fmt.Printf("second counting %v\n", i)
}
