package bridge

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// StartMCPServer expone el Model Context Protocol para Agentes Autónomos.
// Permite a inteligencias artificiales orquestar el nodo local IPv7 de manera determinista.
func StartMCPServer(info *NodeInfo) {
	http.HandleFunc("/mcp/v1/intent", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)

		fmt.Printf("🤖 [MCP Server] Agente Autónomo detectado. Ejecutando intent: %v\n", req["action"])
		// Aquí los agentes pueden modificar topología, pedir estado DHT, o alterar Slices

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"net_state": "stochastic_resonance_active",
			"active_peers": len(info.DHT.Buckets[int(info.DHT.LocalID[0])%160]),
		})
	})

	port := info.APIPort + 100 // MCP usa puerto API + 100
	fmt.Printf("🔌 [MCP Server] Model Context Protocol agentico inicializado en localhost:%d\n", port)
	go http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil)
}
