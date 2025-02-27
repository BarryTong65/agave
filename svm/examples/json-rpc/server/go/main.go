package main

import (
	"fmt"

	"github.com/barry/json-rpc-server/binding"
)

func main() {
	fmt.Println("Simulating transaction...")

	// Call the Rust implementation
	result := binding.SimulateTransaction()

	if result {
		fmt.Printf("%v", result)
	} else {
		fmt.Println("Transaction simulation failed!")
	}
}
