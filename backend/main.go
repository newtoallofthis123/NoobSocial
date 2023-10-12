package main

import "fmt"

func main() {
	fmt.Println("Starting server...")
	api := NewApiServer()
	api.Start()
}
