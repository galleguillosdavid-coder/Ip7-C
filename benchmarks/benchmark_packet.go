//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"time"

	"github.com/galleguillosdavid-coder/Ip7-C/core/protocol"
)

func buildPacket(addr protocol.IPv7Address, subPort uint16, payload []byte) []byte {
	addrWithSP := addr
	addrWithSP.SubPort = subPort
	header := addrWithSP.SerializeHeader()
	sig := protocol.GenerateSignature(payload)
	packet := append(header, sig...)
	packet = append(packet, payload...)
	return packet
}

func main() {
	addr := protocol.NewIPv7(56, 1, 100)
	payload := []byte("test payload")

	start := time.Now()
	for i := 0; i < 100; i++ {
		packet := buildPacket(addr, 0, payload)
		if len(packet) == 0 {
			fmt.Println("Build failed")
			return
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("100 packets built took %v\n", elapsed)
	fmt.Printf("Average per packet: %v\n", elapsed/100)
}
