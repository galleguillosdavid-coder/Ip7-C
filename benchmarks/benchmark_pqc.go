//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"time"

	"github.com/galleguillosdavid-coder/Ip7-C/core/protocol"
)

func main() {
	data := []byte("test data for signature")

	start := time.Now()
	for i := 0; i < 100; i++ {
		sig := protocol.GenerateSignature(data)
		if sig == nil {
			fmt.Println("Signature failed")
			return
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("100 signatures took %v\n", elapsed)
	fmt.Printf("Average per signature: %v\n", elapsed/100)
}
