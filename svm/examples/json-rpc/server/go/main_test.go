package main

import (
	"testing"

	"github.com/barry/json-rpc-server/binding"
)

func BenchmarkSimulateTransaction(b *testing.B) {
	// b.N will be set automatically by the testing framework
	// It will run the loop enough times to get a stable measurement
	for i := 0; i < b.N; i++ {
		binding.SimulateTransaction()
	}
}
