package main

import "fmt"

func init() {
	fmt.Println("Init func")
}

func main() {
	for i := 0; i < 10; i++ {
		fmt.Print("Hello", i, "World")
	}
}
