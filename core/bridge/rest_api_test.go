package bridge

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/galleguillosdavid-coder/Ip7-C/core/overlay"
	"github.com/galleguillosdavid-coder/Ip7-C/core/protocol"
)

func TestHandleConfig(t *testing.T) {
	// Create a mock tunnel with a random high port
	localNode := &protocol.Node{
		Address: protocol.NewIPv7WithSubPort(56, 1, 100, 0),
		Latency: 10.0,
	}
	tunnel, err := overlay.NewTunnel(localNode, 17778, "", 17778, false, "auto")
	if err != nil {
		t.Fatalf("Failed to create tunnel: %v", err)
	}
	defer tunnel.Conn.Close()

	info := &NodeInfo{
		Tunnel: tunnel,
	}

	// Test GET
	req := httptest.NewRequest(http.MethodGet, "/config", nil)
	w := httptest.NewRecorder()
	handleConfig(w, req, info)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if resp["pqc_mode"] != "auto" {
		t.Errorf("Expected pqc_mode 'auto', got %v", resp["pqc_mode"])
	}

	// Test POST
	body := map[string]interface{}{
		"pqc_mode": "off",
	}
	bodyBytes, _ := json.Marshal(body)
	req = httptest.NewRequest(http.MethodPost, "/config", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handleConfig(w, req, info)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check updated values
	if info.Tunnel.GetPQCMode() != "off" {
		t.Errorf("Expected pqc_mode 'off', got %s", info.Tunnel.GetPQCMode())
	}
}
