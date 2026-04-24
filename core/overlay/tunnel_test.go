package overlay

import (
	"testing"
	"time"

	"github.com/galleguillosdavid-coder/Ip7-C/core/protocol"
)

func TestShouldAttachPQC(t *testing.T) {
	localNode := &protocol.Node{
		Address: protocol.NewIPv7WithSubPort(56, 1, 100, 0),
		Latency: 10.0,
	}
	tunnel, err := NewTunnel(localNode, 17779, "", 17779, false, "auto")
	if err != nil {
		t.Fatalf("Failed to create tunnel: %v", err)
	}
	defer tunnel.Conn.Close()

	// Test off mode
	tunnel.mu.Lock()
	tunnel.PqcMode = "off"
	tunnel.mu.Unlock()
	if tunnel.shouldAttachPQC(false) {
		t.Error("Expected false for pqc_mode 'off'")
	}

	// Test noPQC true
	tunnel.mu.Lock()
	tunnel.NoPQC = true
	tunnel.PqcMode = "auto"
	tunnel.mu.Unlock()
	if tunnel.shouldAttachPQC(false) {
		t.Error("Expected false for noPQC true")
	}

	// Test on mode
	tunnel.mu.Lock()
	tunnel.NoPQC = false
	tunnel.PqcMode = "on"
	tunnel.mu.Unlock()
	if !tunnel.shouldAttachPQC(false) {
		t.Error("Expected true for pqc_mode 'on'")
	}

	// Test auto mode, important true
	tunnel.mu.Lock()
	tunnel.PqcMode = "auto"
	tunnel.lastPQCAttach = time.Now().Add(-10 * time.Minute) // Old
	tunnel.mu.Unlock()
	if !tunnel.shouldAttachPQC(true) {
		t.Error("Expected true for important=true")
	}

	// Test auto mode, not important, first time
	tunnel.mu.Lock()
	tunnel.lastPQCAttach = time.Time{}
	tunnel.mu.Unlock()
	if !tunnel.shouldAttachPQC(false) {
		t.Error("Expected true for first attach")
	}
}
