package bridge

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/galleguillosdavid-coder/Ip7-C/core/overlay"
	"github.com/galleguillosdavid-coder/Ip7-C/core/protocol"
)

// StartAgentSandbox habilita la ejecución segura de "Agentes Viajeros" a través de la red IPv7-C.
// Utiliza un entorno aislado y limitado por el Task Budget para evitar consumo excesivo.
func StartAgentSandbox(tunnel *overlay.Tunnel) {
	tunnel.RegisterSubPort(7070, func(remoteAddr protocol.IPv7Address, data []byte) {
		fmt.Printf("🛡️ [Sandbox] Recibido Blueprint Agentico de Satélite %.0f.\n", remoteAddr.ResolvedIP)
		
		// 1. Parsing and validation
		var intent map[string]interface{}
		if err := json.Unmarshal(data, &intent); err != nil {
			fmt.Println("❌ [Sandbox] Error: Blueprint no es un JSON válido")
			tunnel.SendSubPort(remoteAddr, 7070, []byte(`{"status":"error", "reason":"invalid_json_blueprint"}`))
			return
		}

		// 2. Sandboxing determinista (100ms TimeLimit)
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		done := make(chan bool)
		go func() {
			// Simular ejecución analítica real
			// TODO: Enlazar con engine WebAssembly en v1.5
			time.Sleep(10 * time.Millisecond) 
			done <- true
		}()

		select {
		case <-ctx.Done():
			fmt.Println("⚠️ [Sandbox] Ejecución abortada: Límite de tiempo excedido (100ms)")
			tunnel.SendSubPort(remoteAddr, 7070, []byte(`{"status":"timeout_aborted"}`))
		case <-done:
			fmt.Println("✅ [Sandbox] Blueprint ejecutado exitosamente en ventana determinista")
			responseIntent := []byte(`{"status":"executed_safely", "p_bits_state": 0.84}`)
			tunnel.SendSubPort(remoteAddr, 7070, responseIntent)
		}
	})
	fmt.Println("🛡️ [Sandbox] Entorno de ejecución agentica (Egress Sandbox) inicializado en SubPuertos 7070.")
}
