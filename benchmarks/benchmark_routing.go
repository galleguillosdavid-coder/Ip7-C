//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"time"

	"github.com/galleguillosdavid-coder/Ip7-C/core/protocol"
)

func main() {
	node := &protocol.Node{
		Name:    "Test",
		Address: protocol.NewIPv7(56, 1, 100),
		Latency: 10.0,
	}

	// Add some neighbors
	for i := 0; i < 5; i++ {
		neighbor := &protocol.Node{
			Name:    fmt.Sprintf("Neighbor%d", i),
			Address: protocol.NewIPv7(56, 1, float64(101+i)),
			Latency: float64(5 + i),
		}
		node.Neighbors = append(node.Neighbors, neighbor)
	}

	start := time.Now()
	for i := 0; i < 1000; i++ {
		hop := node.NextHop()
		if hop == nil {
			fmt.Println("No hop")
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("1000 NextHop calls took %v\n", elapsed)
	fmt.Printf("Average per call: %v\n", elapsed/1000)
}
